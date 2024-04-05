package disbursement

import (
	"context"
	"disbursement/domain/common/response"
)

type DisbursementService interface {
	CreateFundTransferTransaction(ctx context.Context, req []DisbursementRequest) (err error)
}

type DisbursementRepository interface {
	InsertDisbursement(ctx context.Context, disbursement Disbursement) (err error)
	GetDisbursementsByGroupId(ctx context.Context, groupId string) (res []Disbursement, err error)
	GetDisbursementsById(ctx context.Context, id string) (res Disbursement, err error)
}

type DisbursementUsecase interface {
	CreateDisbursements(ctx context.Context, disbursementRequest []DisbursementRequest) (res *response.Response[DisbursementResponse], err error)
	GetDisbursementsByGroupId(ctx context.Context, groupId string) (res *response.Response[[]Disbursement], err error)
	GetDisbursementById(ctx context.Context, id string) (res *response.Response[Disbursement], err error)
}

type Disbursement struct {
	ID            string `json:"id"`
	GroupId       string `json:"group_id"` // intended for bulk disbursement
	Amount        int64  `json:"amount"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type DisbursementRequest struct {
	Amount        int64  `json:"amount"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

type DisbursementResponse struct {
	ID      *string `json:"id,omitempty"`
	GroupId *string `json:"group_id,omitempty"`
}
