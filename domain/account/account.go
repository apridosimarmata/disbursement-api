package account

import (
	"context"
	"disbursement/domain/common/response"
)

type AccountService interface {
	GetAccountByNumber(ctx context.Context, number string) (res *response.Response[Account], err error) // the bank should not expose any account that is not available for transaction (blocked account -> not found)
}

type AccountUsecase interface {
	GetAccountByNumber(ctx context.Context, number string) (res *response.Response[Account], err error)
}

type Account struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}
