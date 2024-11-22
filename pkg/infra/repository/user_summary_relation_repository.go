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
	db *gorm.DB
}

func NewUserSummaryRelationRepository(db *gorm.DB) repository.IUserSummaryRelation {
	return &userSummaryRelationRepository{
		db: db,
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

func (repo *userSummaryRelationRepository) FindOtherPlatformVc(ctx context.Context, userID, platformNumber uint) (*entity.UserSummaryRelation, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var userSummaryRelation entity.UserSummaryRelation
	if err := tx.Where("user_id = ? AND platform_number != ?", userID, platformNumber).
		First(&userSummaryRelation).Error; err != nil {
		return nil, err
	}

	return &userSummaryRelation, nil
}

func (repo *userSummaryRelationRepository) CreateOrUpdate(ctx context.Context, relation *entity.UserSummaryRelation) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	if relation.ID != 0 {
		t := entity.UserSummaryRelation{Model: gorm.Model{ID: relation.ID}}
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t)
	}

	return tx.Omit(clause.Associations).Save(relation).Error
}
