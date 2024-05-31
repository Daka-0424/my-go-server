package repository

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSummaryRelationRepository struct {
	db     *gorm.DB
	fields []string
}

func NewUserSummaryRelationRepository(db *gorm.DB) repository.IUserSummaryRelation {
	return &userSummaryRelationRepository{
		db:     db,
		fields: entity.GetEntityFields(entity.UserSummaryRelation{}),
	}
}

func (repo *userSummaryRelationRepository) FindByUserID(ctx context.Context, userID uint) (*entity.UserSummaryRelation, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var userSummaryRelation entity.UserSummaryRelation
	if err := tx.Where("user_id = ?", userID).First(&userSummaryRelation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &userSummaryRelation, nil
		}
		return nil, err
	}

	return &userSummaryRelation, nil
}

func (repo *userSummaryRelationRepository) CreateOrUpdate(ctx context.Context, entity *entity.UserSummaryRelation) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.fields),
	}).Create(entity).Error
}
