package repository

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userResourceRepository[T entity.IUserResourceType] struct {
	db     *gorm.DB
	fields []string
}

func NewUserResourceRepository[T entity.IUserResourceType](db *gorm.DB) repository.IUserResource[T] {
	return &userResourceRepository[T]{
		db:     db,
		fields: entity.GetEntityFields(T{}),
	}
}

func (repo *userResourceRepository[T]) GetByID(ctx context.Context, ID uint) (*T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var entity T
	if err := tx.WithContext(ctx).Model(&entity).Where("id = ?", ID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}

func (repo *userResourceRepository[T]) GetByIDs(ctx context.Context, IDs []uint) ([]T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var entities []T
	if err := tx.WithContext(ctx).Model(&entities).Where("id IN ?", IDs).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (repo *userResourceRepository[T]) GetByUserID(ctx context.Context, userID uint) ([]T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var entities []T
	if err := tx.WithContext(ctx).Model(&entities).Where("user_id = ?", userID).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (repo *userResourceRepository[T]) CreateOrUpdate(ctx context.Context, resource *T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.fields),
	}).Create(resource).Error
}

func (repo *userResourceRepository[T]) BulkCreate(ctx context.Context, resources []T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Omit(clause.Associations).Create(resources).Error
}

func (repo *userResourceRepository[T]) BulkUpdate(ctx context.Context, resources []T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(repo.fields),
	}).Create(resources).Error
}
