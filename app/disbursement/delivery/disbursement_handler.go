package delivery

import (
	"disbursement/domain"
	"disbursement/domain/common/response"
	"disbursement/domain/disbursement"
	"disbursement/infrastructure"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type disbursementHandler struct {
	disbursementUsecase   disbursement.DisbursementUsecase
	httpRequestMiddleware infrastructure.HttpRequestMiddleware
}

func SetDisbursementHandler(router *chi.Mux, usecases domain.Usecases, httpRequestMiddleware infrastructure.HttpRequestMiddleware) {
	disbursementHandler := disbursementHandler{
		disbursementUsecase:   usecases.DisbursementUsecase,
		httpRequestMiddleware: httpRequestMiddleware,
	}

	router.Route("/api/v1/disbursements", func(r chi.Router) {
		// POST
		r.Post("/", disbursementHandler.CreateDisburesements)

		// GET
		r.Get("/{id}", disbursementHandler.GetDisbursementById)
	})

	router.Route("/api/v1/disbursements/callback", func(r chi.Router) {
		r.Use(disbursementHandler.httpRequestMiddleware.AuthorizeCallbackRequestMiddleware)
		// POST
		r.Post("/", disbursementHandler.HandleDisbursementCallback)
	})
}

func (disbursementHandler *disbursementHandler) GetDisbursementById(w http.ResponseWriter, r *http.Request) {
	disbursementId := chi.URLParam(r, "id")

	result, err := disbursementHandler.disbursementUsecase.GetDisbursementById(r.Context(), disbursementId)
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result.Success(response.STATUS_SUCCESS, *result.Data)
	result.WriteResponse(w)
}

func (disbursementHandler *disbursementHandler) HandleDisbursementCallback(w http.ResponseWriter, r *http.Request) {
	req := disbursement.DisbursementCallbackRequest{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	if err := req.Validate(); err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result, err := disbursementHandler.disbursementUsecase.HandleDisbursementCallback(r.Context(), req)
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result.Success(response.STATUS_SUCCESS, *result.Data)
	result.WriteResponse(w)
}

func (disbursementHandler *disbursementHandler) CreateDisburesements(w http.ResponseWriter, r *http.Request) {
	req := disbursement.BulkDisbursementRequest{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result, err := disbursementHandler.disbursementUsecase.CreateDisbursements(r.Context(), req.Data)
	if err != nil {
		errResp := &response.Response[response.Error]{
			Data: &response.Error{
				Error: err.Error(),
			},
		}
		errResp.Error(err.Error())
		errResp.WriteResponse(w)
		return
	}

	result.Success(response.STATUS_SUCCESS, *result.Data)
	result.WriteResponse(w)
}
