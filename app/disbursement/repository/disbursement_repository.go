package disbursement

import (
	"context"
	"database/sql"
	"disbursement/domain/disbursement"

	sq "github.com/Masterminds/squirrel"

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

func (disbursementRepository *disbursementRepository) InsertDisbursements(ctx context.Context, disbursements []disbursement.Disbursement) (err error) {
	err = disbursementRepository.db.WithContext(ctx).Table("tx_disbursements").Create(disbursements).Error
	if err != nil {
		return err
	}

	return nil
}

func (disbursementRepository *disbursementRepository) UpdateDisbursementStatus(ctx context.Context, status string, disbursementIds []string) (err error) {
	err = disbursementRepository.db.WithContext(ctx).Table("tx_disbursements").Where("id IN ?", disbursementIds).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}

func (disbursementRepository *disbursementRepository) GetDisbursementsByIds(ctx context.Context, ids []string) (res []disbursement.Disbursement, err error) {
	builder := sq.Select("*").From("tx_disbursements").Where("id IN ?", ids)
	qry, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = disbursementRepository.db.WithContext(ctx).Raw(qry, args...).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
