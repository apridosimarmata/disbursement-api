package disbursement

import (
	"context"
	"disbursement/domain"
	"disbursement/domain/common/response"
	disbursement "disbursement/domain/disbursement"
	"disbursement/infrastructure"
	"errors"
	"fmt"
)

type disbursementUsecase struct {
	disbursementService    disbursement.DisbursementService
	disbursementRepository disbursement.DisbursementRepository
}

func NewDisbursementUsecase(
	repositories domain.Repositories,
	services domain.Services,
) disbursement.DisbursementUsecase {
	return &disbursementUsecase{
		disbursementService:    services.DisbursementService,
		disbursementRepository: repositories.DisbursementRepository,
	}
}

func (disbursementUsecase *disbursementUsecase) CreateDisbursements(ctx context.Context, disbursementRequest []disbursement.DisbursementRequest) (res *response.Response[disbursement.DisbursementResponse], err error) {
	// save to db

	// hit service

	return res, nil
}

func (disbursementUsecase *disbursementUsecase) GetDisbursementsByGroupId(ctx context.Context, groupId string) (res *response.Response[[]disbursement.Disbursement], err error) {
	disbursements, err := disbursementUsecase.disbursementRepository.GetDisbursementsByGroupId(ctx, groupId)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.disbursementRepository.GetDisbursementsByGroupId @ disbursementUsecase.GetDisbursementsByGroupId", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	res = &response.Response[[]disbursement.Disbursement]{
		Data: &disbursements,
	}
	return res, nil
}

func (disbursementUsecase *disbursementUsecase) GetDisbursementById(ctx context.Context, id string) (res *response.Response[disbursement.Disbursement], err error) {
	_disbursement, err := disbursementUsecase.disbursementRepository.GetDisbursementsById(ctx, id)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.disbursementRepository.GetDisbursementsById @ disbursementUsecase.GetDisbursementById", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	res = &response.Response[disbursement.Disbursement]{
		Data: &_disbursement,
	}
	return res, nil

}
