package disbursement

import (
	"context"
	"disbursement/domain"
	"disbursement/domain/account"
	"disbursement/domain/common/response"
	disbursement "disbursement/domain/disbursement"
	"disbursement/infrastructure"
	"errors"
	"fmt"

	uuid "github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type disbursementUsecase struct {
	disbursementService    disbursement.DisbursementService
	disbursementRepository disbursement.DisbursementRepository
	accountService         account.AccountService
}

func NewDisbursementUsecase(
	repositories domain.Repositories,
	services domain.Services,
) disbursement.DisbursementUsecase {
	return &disbursementUsecase{
		disbursementService:    services.DisbursementService,
		disbursementRepository: repositories.DisbursementRepository,
		accountService:         services.AccountService,
	}
}

func (disbursementUsecase *disbursementUsecase) sendDisbursementRequests(ctx context.Context, disbursements []disbursement.Disbursement) (err error) {
	// the bank account disburse API only accept single disbursement per request
	// do the disbursement concurrently
	var group errgroup.Group
	failedDisbursmentIdsChannel := make(chan string, len(disbursements))
	group.SetLimit(len(disbursements))

	for _, disb := range disbursements {
		disbursementCopy := disb
		group.Go(func() error {
			err := disbursementUsecase.disbursementService.DisburseFund(ctx, disbursement.DisburseFundRequest{
				DisbursementID: disbursementCopy.ID,
				Amount:         disbursementCopy.Amount,
				AccountNumber:  disbursementCopy.AccountNumber,
			})
			if err != nil {
				failedDisbursmentIdsChannel <- disbursementCopy.ID
			}

			return nil
		})
	}

	group.Wait()
	close(failedDisbursmentIdsChannel)

	// extracting failed disbursement ids from channel
	// immidiately update status to fail
	// note: bank callback also can fail the disbursement
	failedDisbursementIds := []string{}
	for disbursementId := range failedDisbursmentIdsChannel {
		failedDisbursementIds = append(failedDisbursementIds, disbursementId)
	}

	go catchFailingBankDisbursementRequest(ctx, failedDisbursementIds) // do the update asynchronously, avoid late response to client

	return nil
}

func catchFailingBankDisbursementRequest(ctx context.Context, failingDisbursementIds []string) {

}

func (disbursementUsecase *disbursementUsecase) validateAccounts(ctx context.Context, accountNumbers []string) (validAccount map[string]string, err error) {
	// the bank account validation API only accept single account number per request
	// do the validation concurrently
	var group errgroup.Group
	validAccountChannel := make(chan map[string]string, len(accountNumbers))
	group.SetLimit(len(accountNumbers))

	for _, accountNumber := range accountNumbers {
		// param should be copied for go with version earlier than v1.22
		// https://medium.com/@isimarmata09/how-i-use-goroutine-in-the-wrong-way-8cb4e8282efb
		accountNumberCopy := accountNumber
		group.Go(func() error {
			account, err := disbursementUsecase.accountService.GetAccountByNumber(ctx, accountNumberCopy)
			if err != nil && err.Error() != response.ERROR_NOT_FOUND { // if error occured other than not found, catch error
				return err
			}
			if err != nil && err.Error() == response.ERROR_NOT_FOUND {
				return nil
			}

			validAccountChannel <- map[string]string{
				account.Number: account.Name,
			}

			return nil
		})
	}

	if err = group.Wait(); err != nil {
		return nil, err
	}

	close(validAccountChannel)

	// extracting valid account numbers from channel
	validAccount = make(map[string]string)
	for account := range validAccountChannel {
		for accountNumber, accountName := range account {
			validAccount[accountNumber] = accountName
		}

	}

	return validAccount, err
}

// designed to accept 1 to many disbursement request at one time
func (disbursementUsecase *disbursementUsecase) CreateDisbursements(ctx context.Context, disbursementRequest []disbursement.DisbursementRequest) (res *response.Response[disbursement.DisbursementResponse], err error) {
	var singleDisbursementId *string
	disbursements := []disbursement.Disbursement{}

	accountNumbers := []string{}
	for _, disbursementReq := range disbursementRequest {
		accountNumbers = append(accountNumbers, disbursementReq.AccountNumber)
	}

	validAccountMap, err := disbursementUsecase.validateAccounts(ctx, accountNumbers)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.validateAccounts @ disbursementUsecase.CreateDisbursements", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	groupId, err := uuid.NewV6()
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - uuid.NewV6 @ disbursementUsecase.CreateDisbursements", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}
	for _, disb := range disbursementRequest {
		disbursementId, err := uuid.NewV6()
		disbursementIdString := disbursementId.String()
		singleDisbursementId = &disbursementIdString
		if err != nil {
			infrastructure.Log(fmt.Sprintf("%s - uuid.NewV6 @ disbursementUsecase.CreateDisbursements", err.Error()))
			return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
		}

		if _, validAccount := validAccountMap[disb.AccountNumber]; validAccount {
			disbursements = append(disbursements, disbursement.Disbursement{
				ID:            disbursementId.String(),
				GroupId:       groupId.String(),
				Amount:        disb.Amount,
				AccountNumber: disb.AccountNumber,
				AccountName:   validAccountMap[disb.AccountNumber],
				Status:        disbursement.DISBURSEMENT_STATUS_PENDING,
			})
		} else {
			disbursements = append(disbursements, disbursement.Disbursement{
				ID:            disbursementId.String(),
				GroupId:       groupId.String(),
				Amount:        disb.Amount,
				AccountNumber: disb.AccountNumber,
				AccountName:   "-",
				Status:        disbursement.DISBURSEMENT_STATUS_FAILED,
				Message:       disbursement.DISBURSEMENT_INVALID_ACCOUNT,
			})
		}
	}

	// write disbursements to database
	err = disbursementUsecase.disbursementRepository.InsertDisbursements(ctx, disbursements)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.disbursementRepository.InsertDisbursements @ disbursementUsecase.CreateDisbursements", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	// send request to bank
	err = disbursementUsecase.sendDisbursementRequests(ctx, disbursements)
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.sendDisbursementRequests @ disbursementUsecase.CreateDisbursements", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	result := disbursement.DisbursementResponse{}
	if len(disbursements) < 2 {
		result.ID = singleDisbursementId
	} else {
		groupIdString := groupId.String()
		result.GroupId = &groupIdString
	}

	return &response.Response[disbursement.DisbursementResponse]{
		Data: &result,
	}, nil
}

func (disbursementUsecase *disbursementUsecase) HandleDisbursementCallback(ctx context.Context, disbursementCallbackRequest disbursement.DisbursementCallbackRequest) (res *response.Response[string], err error) {
	// assumption:
	// return err to demand a retry from bank to send the callback again (until success)

	disbursementEntities, err := disbursementUsecase.disbursementRepository.GetDisbursementsByIds(ctx, []string{disbursementCallbackRequest.DisbursementID})
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.disbursementRepository.GetDisbursementsById @ disbursementUsecase.HandleDisbursementCallback", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	if len(disbursementEntities) == 0 {
		return nil, errors.New(response.ERROR_DISBURSEMENT_NOT_FOUND)
	}

	disbursementEntity := disbursementEntities[0]

	if disbursementEntity.Status != disbursement.DISBURSEMENT_STATUS_PENDING {
		return nil, errors.New(response.ERROR_DISBURSEMENT_STATUS_ALREADY_UPDATED_BEFORE)
	}

	err = disbursementUsecase.disbursementRepository.UpdateDisbursementStatus(ctx, disbursementCallbackRequest.Status, []string{disbursementEntity.ID})
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.disbursementRepository.UpdateDisbursement @ disbursementUsecase.HandleDisbursementCallback", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	var success = "success"
	return &response.Response[string]{
		Data: &success,
	}, nil
}

func (disbursementUsecase *disbursementUsecase) GetDisbursementById(ctx context.Context, id string) (res *response.Response[disbursement.Disbursement], err error) {
	_disbursement, err := disbursementUsecase.disbursementRepository.GetDisbursementsByIds(ctx, []string{id})
	if err != nil {
		infrastructure.Log(fmt.Sprintf("%s - disbursementUsecase.disbursementRepository.GetDisbursementsById @ disbursementUsecase.GetDisbursementById", err.Error()))
		return nil, errors.New(response.ERROR_INTERNAL_SERVER_ERROR)
	}

	if len(_disbursement) == 0 {
		return nil, errors.New(response.ERROR_DISBURSEMENT_NOT_FOUND)
	}

	res = &response.Response[disbursement.Disbursement]{
		Data: &_disbursement[0],
	}
	return res, nil

}
