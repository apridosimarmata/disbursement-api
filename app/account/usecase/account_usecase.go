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
	result, err := accountUsecase.accountService.GetAccountByNumber(ctx, number)
	if err != nil && err.Error() != response.ERROR_NOT_FOUND {
		infrastructure.Log(fmt.Sprintf("%s - accountUsecase.accountService.GetAccountByNumber @ accountUsecase.GetAccountByNumber", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	if err != nil && err.Error() == response.ERROR_NOT_FOUND {
		infrastructure.Log(fmt.Sprintf("%s - accountUsecase.accountService.GetAccountByNumber @ accountUsecase.GetAccountByNumber", err.Error()))
		return nil, errors.New(response.ERROR_ACCOUNT_NOT_FOUND)
	}

	return &response.Response[account.Account]{
		Data: result,
	}, nil
}
