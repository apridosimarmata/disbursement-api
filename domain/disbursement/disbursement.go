package disbursement

import (
	"context"
	"disbursement/domain/common/response"
	"errors"
)

const (
	DISBURSEMENT_STATUS_PENDING = "pending"
	DISBURSEMENT_STATUS_SUCCESS = "success"
	DISBURSEMENT_STATUS_FAILED  = "failed"

	DISBURSEMENT_INVALID_ACCOUNT = "invalid account"
)

var disbursementStatusMap = map[string]struct{}{
	DISBURSEMENT_STATUS_SUCCESS: {},
	DISBURSEMENT_STATUS_PENDING: {},
	DISBURSEMENT_STATUS_FAILED:  {},
}

type DisbursementService interface {
	DisburseFund(ctx context.Context, req DisburseFundRequest) (err error)
}

type DisbursementRepository interface {
	InsertDisbursements(ctx context.Context, disbursements []Disbursement) (err error)
	UpdateDisbursementStatus(ctx context.Context, status string, disbursementIds []string) (err error)
	GetDisbursementsByIds(ctx context.Context, ids []string) (res []Disbursement, err error)
}

type DisbursementUsecase interface {
	CreateDisbursements(ctx context.Context, disbursementRequest []DisbursementRequest) (res *response.Response[DisbursementResponse], err error)
	HandleDisbursementCallback(ctx context.Context, disbursementCallbackRequest DisbursementCallbackRequest) (res *response.Response[string], err error)
	// GetDisbursementsByGroupId(ctx context.Context, groupId string) (res *response.Response[[]Disbursement], err error)
	GetDisbursementById(ctx context.Context, id string) (res *response.Response[Disbursement], err error)
}

type Disbursement struct {
	ID            string `json:"id"`
	GroupId       string `json:"group_id"` // intended for to group a bulk disbursement
	Amount        int64  `json:"amount"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type DisbursementRequest struct {
	Amount        int64  `json:"amount"`
	AccountNumber string `json:"account_number"`
}

type BulkDisbursementRequest struct {
	Data []DisbursementRequest `json:"data"`
}

type DisburseFundRequest struct {
	DisbursementID string `json:"disbursement_id"`
	Amount         int64  `json:"amount"`
	AccountNumber  string `json:"account_number"`
}

type DisbursementResponse struct {
	ID      *string `json:"id,omitempty"`
	GroupId *string `json:"group_id,omitempty"`
}

type DisbursementCallbackRequest struct {
	DisbursementID string `json:"disbursement_id"`
	Status         string `json:"status"`
}

func (payload *DisbursementCallbackRequest) Validate() (err error) {
	if _, valid := disbursementStatusMap[payload.Status]; !valid {
		return errors.New(response.ERROR_UNRECOGNIZED_DISBURSEMENT_STATUS)
	}

	return nil
}
