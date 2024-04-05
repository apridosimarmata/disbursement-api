package account

import (
	"context"
	"disbursement/domain"
	"disbursement/domain/account"
	"disbursement/domain/common/response"
	"disbursement/infrastructure"
	"errors"
	"fmt"
)

type accountUsecase struct {
	accountService account.AccountService
}

func NewAccountUsecase(services domain.Services) account.AccountUsecase {
	return &accountUsecase{
		accountService: services.AccountService,
	}
}

func (accountUsecase *accountUsecase) GetAccountByNumber(ctx context.Context, number string) (res *response.Response[account.Account], err error) {
	account, err := accountUsecase.accountService.GetAccountByNumber(ctx, number)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - accountUsecase.accountService.GetAccountByNumber @ accountUsecase.GetAccountByNumber", err.Error()))
		return nil, errors.New(response.ERROR_ACCOUNT_NOT_FOUND)
	}

	return account, nil
}
