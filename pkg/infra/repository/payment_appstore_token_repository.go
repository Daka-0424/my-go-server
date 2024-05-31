package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentAppstoreTokenRepository struct {
	db     *gorm.DB
	fields []string
}

func NewPaymentAppstoreTokenRepository(db *gorm.DB) repository.IPaymentAppstoreToken {
	return &paymentAppstoreTokenRepository{
		db:     db,
		fields: entity.GetEntityFields(entity.PaymentAppstoreToken{}),
	}
}

func (repo *paymentAppstoreTokenRepository) CreateOrUpdate(ctx context.Context, appstoreToken *entity.PaymentAppstoreToken) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.fields),
	}).Create(appstoreToken).Error
}

func (r *paymentAppstoreTokenRepository) ExistsPaymentAppstoreToken(ctx context.Context, transactionID string) (bool, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		return false, repository.ErrTx
	}

	var count int64
	if err := tx.Model(&entity.PaymentAppstoreToken{}).Where("transaction_id = ?", transactionID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
