package disbursement

import (
	"context"
	"disbursement/domain/common/response"
	"disbursement/domain/disbursement"
	"disbursement/infrastructure"
	"fmt"
)

type disbursementService struct {
	dummyBankServiceClient infrastructure.HTTPClient
}

func NewDisbursementService(dummyBankServiceClient infrastructure.HTTPClient) disbursement.DisbursementService {
	return &disbursementService{
		dummyBankServiceClient: dummyBankServiceClient,
	}
}

func (disbursementService *disbursementService) CreateFundTransferTransaction(ctx context.Context, req []disbursement.DisbursementRequest) (err error) {
	var response response.Response[string]

	err = disbursementService.dummyBankServiceClient.Post(ctx, fmt.Sprintf("/accounts/%s", req), nil, response)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementService.httpClient.Post @ accountService.CreateFundTransferTransaction", err.Error()))
		return err
	}

	if response.StatusCode != 200 {
		infrastructure.Log(fmt.Sprintf("%s - @ accountService.CreateFundTransferTransaction", response.Status))
		return
	}

	return nil
}
