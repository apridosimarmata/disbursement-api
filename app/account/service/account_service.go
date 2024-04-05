package account

import (
	"context"
	"disbursement/domain/account"
	"disbursement/domain/common/response"
	"disbursement/infrastructure"
	"fmt"
)

type accountService struct {
	dummyBankServiceClient infrastructure.HTTPClient
}

func NewAccountService(dummyBankServiceClient infrastructure.HTTPClient) account.AccountService {
	return &accountService{
		dummyBankServiceClient: dummyBankServiceClient,
	}
}

func (accountService *accountService) GetAccountByNumber(ctx context.Context, number string) (res *account.Account, err error) {
	response := &response.Response[account.Account]{}
	err = accountService.dummyBankServiceClient.Get(ctx, fmt.Sprintf("/accounts/%s", number), nil, response)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - accountService.httpClient.Get @ accountService.GetAccountByNumber", err.Error()))
		return nil, err
	}

	return response.Data, nil
}
