package domain

import (
	"disbursement/domain/account"
	"disbursement/domain/disbursement"
)

type Repositories struct {
	DisbursementRepository disbursement.DisbursementRepository
}

type Services struct {
	AccountService      account.AccountService
	DisbursementService disbursement.DisbursementService
}

type Usecases struct {
	AccountUsecase      account.AccountUsecase
	DisbursementUsecase disbursement.DisbursementUsecase
}
