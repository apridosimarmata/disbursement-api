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

func (disbursementService *disbursementService) DisburseFund(ctx context.Context, req disbursement.DisburseFundRequest) (err error) {
	var response response.Response[string]

	err = disbursementService.dummyBankServiceClient.Post(ctx, "/disbursements", req, response)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementService.httpClient.Post @ accountService.DisburseFund", err.Error()))
		return err
	}

	return nil
}
