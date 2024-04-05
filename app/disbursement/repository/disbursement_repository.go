package disbursement

import (
	"context"
	"disbursement/domain/disbursement"

	"gorm.io/gorm"
)

type disbursementRepository struct {
	db *gorm.DB
}

func NewDisbursementRepository(db *gorm.DB) disbursement.DisbursementRepository {
	return &disbursementRepository{
		db: db,
	}
}

func (disbursementRepository *disbursementRepository) InsertDisbursement(ctx context.Context, disbursement disbursement.Disbursement) (err error) {
	return nil
}

func (disbursementRepository *disbursementRepository) GetDisbursementsByGroupId(ctx context.Context, groupId string) (res []disbursement.Disbursement, err error) {

	return res, nil
}

func (disbursementRepository *disbursementRepository) GetDisbursementsById(ctx context.Context, id string) (res disbursement.Disbursement, err error) {

	return res, nil
}
