package repository

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
)

type seedRepository[T entity.ISeedType] struct {
	db *gorm.DB
}

func NewSeedRepository[T entity.ISeedType](db *gorm.DB) repository.ISeed[T] {
	return &seedRepository[T]{
		db: db,
	}
}

func (r *seedRepository[T]) GetByID(ctx context.Context, ID uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Where("id = ?", ID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}

func (r *seedRepository[T]) GetAll(ctx context.Context, preloads ...string) ([]T, error) {
	db := r.db
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var entities []T
	if err := db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *seedRepository[T]) Get(ctx context.Context, param T) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Where(&param).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}

func (r *seedRepository[T]) Where(ctx context.Context, param T) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Where(&param).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *seedRepository[T]) GetByIDs(ctx context.Context, IDs []uint) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Model(&entities).Where(IDs).Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return entities, nil
}

func (r *seedRepository[T]) GetByMap(ctx context.Context, param map[string]interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Where(param).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}
