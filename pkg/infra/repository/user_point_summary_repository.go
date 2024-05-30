package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type userPointSummaryRepository struct {
	db     *gorm.DB
	fields []string
}

func NewUserPointSummaryRepository(db *gorm.DB) repository.IUserPointSummary {
	return &userPointSummaryRepository{
		db:     db,
		fields: entity.GetEntityFields(entity.UserPointSummary{}),
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

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.fields),
	}).Create(pointSummary).Error
}

func (repo *userPointSummaryRepository) BulkUpdate(ctx context.Context, points []entity.UserPointSummary) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.fields),
	}).Create(points).Error
}
