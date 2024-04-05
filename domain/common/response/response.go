package response

import (
	"encoding/json"
	"net/http"
)

const (
	STATUS_SUCCESS = "success"
	STATUS_FAIL    = "fail"
	STATUS_ERROR   = "error"

	ERROR_ACCOUNT_NOT_FOUND     = "account not found"
	ERROR_UNAUTHORIZED          = "unauthorized"
	ERROR_INTERNAL_SERVER_ERROR = "an error occured during processing your request, try again later"

	ERROR_DISBURSEMENT_NOT_FOUND                     = "no such disbursement"
	ERROR_UNRECOGNIZED_DISBURSEMENT_STATUS           = "unrecognized disbursement status"
	ERROR_DISBURSEMENT_STATUS_ALREADY_UPDATED_BEFORE = "could not update disbursement status again"

	// for http client interacting with external service
	ERROR_NOT_FOUND = "not found"
)

var (
	userErrors = map[string]struct{}{
		ERROR_ACCOUNT_NOT_FOUND:                          {},
		ERROR_UNAUTHORIZED:                               {},
		ERROR_DISBURSEMENT_NOT_FOUND:                     {},
		ERROR_UNRECOGNIZED_DISBURSEMENT_STATUS:           {},
		ERROR_DISBURSEMENT_STATUS_ALREADY_UPDATED_BEFORE: {},
	}
)

type Error struct {
	Error string `json:"error"`
}

type Response[T any] struct {
	Status     string `json:"status"`
	Data       *T     `json:"data,omitempty"`
	StatusCode int    `json:"-"`
}

func (payload *Response[T]) Success(msg string, data T) {
	payload.Status = STATUS_SUCCESS
	payload.StatusCode = http.StatusOK
}

func (payload *Response[T]) Error(msg string) {
	payload.Status = STATUS_ERROR
	payload.StatusCode = http.StatusInternalServerError

	if _, isUserError := userErrors[msg]; isUserError {
		payload.Status = STATUS_FAIL
		payload.StatusCode = http.StatusBadRequest
	}
}

func (res *Response[T]) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}
