package delivery

import (
	"disbursement/domain"
	"disbursement/domain/account"
	"net/http"

	"disbursement/domain/common/response"

	"github.com/go-chi/chi/v5"
)

type accountHandler struct {
	accountUsecase account.AccountUsecase
}

func SetAccountHandler(router *chi.Mux, usecases domain.Usecases) {
	accountHandler := accountHandler{
		accountUsecase: usecases.AccountUsecase,
	}

	router.Route("/api/v1/accounts", func(r chi.Router) {

		// GET
		r.Get("/{number}", accountHandler.GetAccountByNumber)
	})

}

func (accountHandler *accountHandler) GetAccountByNumber(w http.ResponseWriter, r *http.Request) {
	accountNumber := chi.URLParam(r, "number")

	result, err := accountHandler.accountUsecase.GetAccountByNumber(r.Context(), accountNumber)
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
