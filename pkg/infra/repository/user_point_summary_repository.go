package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type userPointSummaryRepository struct {
	db *gorm.DB
}

func NewUserPointSummaryRepository(db *gorm.DB) repository.IUserPointSummary {
	return &userPointSummaryRepository{
		db: db,
	}
}

func (repo *userPointSummaryRepository) Find(ctx context.Context, userID, platformNumber uint, paidKind int) (*entity.UserPointSummary, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var userPointSummary entity.UserPointSummary
	if err := tx.Where("user_id = ? AND platform_number = ? AND paid_kind = ?", userID, platformNumber, paidKind).First(&userPointSummary).Error; err != nil {
		return nil, err
	}

	return &userPointSummary, nil
}

func (repo *userPointSummaryRepository) FirstOrCreateFreePointSummary(ctx context.Context, userID uint) (*entity.UserPointSummary, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var userPointSummary entity.UserPointSummary
	if err := tx.Where("user_id = ? AND paid_kind = ?", userID, entity.GemKindFree).FirstOrCreate(&userPointSummary).Error; err != nil {
		return nil, err
	}

	return &userPointSummary, nil
}

func (repo *userPointSummaryRepository) Update(ctx context.Context, pointSummary *entity.UserPointSummary) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	t := entity.UserPointSummary{Model: gorm.Model{ID: pointSummary.ID}}
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t)

	return tx.Model(pointSummary).Select("*").Omit(clause.Associations).Updates(pointSummary).Error
}

func (repo *userPointSummaryRepository) BulkUpdate(ctx context.Context, points []entity.UserPointSummary) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	for _, point := range points {
		if err := tx.Model(&entity.UserPointSummary{}).
			Where("id = ?", point.ID).
			Select("*").
			Omit(clause.Associations).
			Updates(point).Error; err != nil {
			return err
		}
	}

	return nil
}
