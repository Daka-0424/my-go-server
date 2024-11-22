package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userResourceRepository[T entity.IUserResourceType] struct {
	db *gorm.DB
}

func NewUserResourceRepository[T entity.IUserResourceType](db *gorm.DB) repository.IUserResource[T] {
	return &userResourceRepository[T]{
		db: db,
	}
}

func (r *userResourceRepository[T]) GetByID(ctx context.Context, ID uint, preloads ...string) (*T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var entity T
	if err := tx.WithContext(ctx).Model(&entity).Where("id = ?", ID).First(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *userResourceRepository[T]) GetByIDs(ctx context.Context, ISs []uint, preloads ...string) ([]T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var entities []T
	if err := tx.WithContext(ctx).Model(&entities).Where("id IN ?", ISs).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *userResourceRepository[T]) GetByUserID(ctx context.Context, userID uint, preloads ...string) ([]T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var entities []T
	if err := tx.WithContext(ctx).Model(&entities).Where("user_id = ?", userID).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *userResourceRepository[T]) Where(ctx context.Context, param T, preloads ...string) ([]T, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var entities []T
	if err := tx.Where(&param).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *userResourceRepository[T]) CreateOrUpdate(ctx context.Context, entity *T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	if (*entity).IsEmpty() {
		t := (*entity)
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t)
	}

	return tx.Omit(clause.Associations).Save(entity).Error
}

func (r *userResourceRepository[T]) BulkCreate(ctx context.Context, entities []T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Omit(clause.Associations).Create(entities).Error
}

func (r *userResourceRepository[T]) BulkUpdate(ctx context.Context, entities []T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	for _, entity := range entities {
		id := entity.GetID()
		if err := tx.Model(&entity).Where("id = ?", id).Select("*").Omit(clause.Associations).Updates(entity).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userResourceRepository[T]) Delete(ctx context.Context, entity *T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Unscoped().Delete(entity).Error
}

func (r *userResourceRepository[T]) BulkDelete(ctx context.Context, entities []T) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Unscoped().Delete(entities).Error
}

func (r *userResourceRepository[T]) DeleteByUserID(ctx context.Context, userID uint) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	var entity T
	return tx.Unscoped().Where("user_id = ?", userID).Delete(&entity).Error
}
