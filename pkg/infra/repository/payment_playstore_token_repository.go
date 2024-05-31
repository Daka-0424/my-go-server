package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type paymentPlaystoreTokenRepository struct {
	db    *gorm.DB
	field []string
}

func NewPaymentPlaystoreTokenRepository(db *gorm.DB) repository.IPaymentPlaystoreToken {
	return &paymentPlaystoreTokenRepository{
		db:    db,
		field: entity.GetEntityFields(entity.PaymentPlaystoreToken{}),
	}
}

func (repo *paymentPlaystoreTokenRepository) CreateOrUpdate(ctx context.Context, playstoreToken *entity.PaymentPlaystoreToken) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.field),
	}).Create(playstoreToken).Error
}

func (r *paymentPlaystoreTokenRepository) ExistsPaymentPlaystoreToken(ctx context.Context, orderID string) (bool, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		return false, repository.ErrTx
	}

	var count int64
	if err := tx.Model(&entity.PaymentPlaystoreToken{}).Where("order_id = ?", orderID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
