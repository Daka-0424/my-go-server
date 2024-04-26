package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
)

type seedRepository[T entity.SeedType] struct {
	db *gorm.DB
}

func NewSeedRepository[T entity.SeedType](db *gorm.DB) repository.Seed[T] {
	return &seedRepository[T]{db: db}
}

func (repo *seedRepository[T]) GetAll(ctx context.Context, preloads ...string) ([]T, error) {
	db := repo.db
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var entitys []T
	if err := db.WithContext(ctx).Find(&entitys).Error; err != nil {
		return nil, err
	}

	return entitys, nil
}

func (repo *seedRepository[T]) Where(ctx context.Context, param T) ([]T, error) {
	var entitys []T
	if err := repo.db.WithContext(ctx).Where(&param).Find(&entitys).Error; err != nil {
		return nil, err
	}

	return entitys, nil
}
