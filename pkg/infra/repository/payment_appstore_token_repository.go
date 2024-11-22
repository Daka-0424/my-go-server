package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentAppstoreTokenRepository struct {
	db *gorm.DB
}

func NewPaymentAppstoreTokenRepository(db *gorm.DB) repository.IPaymentAppstoreToken {
	return &paymentAppstoreTokenRepository{
		db: db,
	}
}

func (repo *paymentAppstoreTokenRepository) CreateOrUpdate(ctx context.Context, appstoreToken *entity.PaymentAppstoreToken) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	if appstoreToken.ID != 0 {
		t := entity.PaymentAppstoreToken{Model: gorm.Model{ID: appstoreToken.ID}}
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t)
	}

	return tx.Omit(clause.Associations).Save(appstoreToken).Error
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
