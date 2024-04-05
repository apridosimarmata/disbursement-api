package infrastructure

import (
	"disbursement/domain/common/response"
	"net/http"
)

type HttpRequestMiddleware interface {
	AuthorizeCallbackRequestMiddleware(next http.Handler) http.Handler
}

type httpRequestMiddleware struct {
	apiKey string
}

func NewHttpRequestMiddleware(apiKey string) HttpRequestMiddleware {
	return &httpRequestMiddleware{
		apiKey: apiKey,
	}
}

func (middleware *httpRequestMiddleware) AuthorizeCallbackRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedResp := response.Response[response.Error]{}

		requestApiKey := r.Header.Get("X-API-KEY")
		if requestApiKey != middleware.apiKey {
			unauthorizedResp.Data = &response.Error{
				Error: response.ERROR_UNAUTHORIZED,
			}
			unauthorizedResp.Error(response.ERROR_UNAUTHORIZED)
			unauthorizedResp.WriteResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
