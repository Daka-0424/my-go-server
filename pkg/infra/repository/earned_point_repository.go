package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type earnedPointRepository struct {
	db *gorm.DB
}

func NewEarnedPointRepository(db *gorm.DB) repository.IEarnedPoint {
	return &earnedPointRepository{
		db: db,
	}
}

func (repo *earnedPointRepository) CreateOrUpdate(ctx context.Context, earnedPoint *entity.EarnedPoint) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	if earnedPoint.ID != 0 {
		t := entity.EarnedPoint{Model: gorm.Model{ID: earnedPoint.ID}}
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t)
	}

	return tx.Omit(clause.Associations).Save(earnedPoint).Error
}

func (repo *earnedPointRepository) BulkCreate(ctx context.Context, earnedPoints []entity.EarnedPoint) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Omit(clause.Associations).Create(earnedPoints).Error
}

func (r *earnedPointRepository) GetAll(ctx context.Context, offset int, limit int) ([]entity.EarnedPoint, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	var earned_points []entity.EarnedPoint
	err := tx.Offset(offset).Limit(limit).Find(&earned_points).Error
	if err != nil {
		return nil, err
	}

	return earned_points, nil
}

func (r *earnedPointRepository) GetWhere(ctx context.Context, param entity.EarnedPoint, offset int, limit int) ([]entity.EarnedPoint, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	var earned_points []entity.EarnedPoint
	err := tx.Where(&param).Offset(offset).Limit(limit).Find(&earned_points).Error
	if err != nil {
		return nil, err
	}

	return earned_points, nil
}

func (r *earnedPointRepository) FindByPointSummaryIDs(ctx context.Context, pointSummaryIDs ...uint) ([]entity.EarnedPoint, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	var earned_points []entity.EarnedPoint
	err := tx.Where("user_point_summary_id IN ?", pointSummaryIDs).Where("point_exceeded = ?", false).Order("earned_at asc").Find(&earned_points).Error
	if err != nil {
		return nil, err
	}

	return earned_points, nil
}

func (r *earnedPointRepository) CountAll(ctx context.Context) (int64, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	count := int64(0)
	err := tx.Model(&entity.EarnedPoint{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *earnedPointRepository) CountWhere(ctx context.Context, param entity.EarnedPoint) (int64, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	count := int64(0)
	err := tx.Model(&entity.EarnedPoint{}).Where(&param).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
