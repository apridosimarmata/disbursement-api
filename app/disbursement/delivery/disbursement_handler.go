package delivery

import (
	"disbursement/domain"
	"disbursement/domain/common/response"
	"disbursement/domain/disbursement"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type disbursementHandler struct {
	disbursementUsecase disbursement.DisbursementUsecase
}

func SetDisbursementHandler(router *chi.Mux, usecases domain.Usecases) {
	disbursementHandler := disbursementHandler{
		disbursementUsecase: usecases.DisbursementUsecase,
	}

	router.Route("/api/v1/disbursements", func(r chi.Router) {

		// GET
		r.Post("/", disbursementHandler.CreateDisburesements)
	})

}

func (disbursementHandler *disbursementHandler) CreateDisburesements(w http.ResponseWriter, r *http.Request) {
	req := []disbursement.DisbursementRequest{}
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

	result, err := disbursementHandler.disbursementUsecase.CreateDisbursements(r.Context(), req)
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
