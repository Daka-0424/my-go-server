package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userUniqueResourceRepository[T entity.UserUniqueResourceType] struct {
	db *gorm.DB
}

func NewUserUniqueResourceRepository[T entity.UserUniqueResourceType](
	db *gorm.DB,
) repository.IUserUniqueResource[T] {
	return &userUniqueResourceRepository[T]{
		db: db,
	}
}

func (r *userUniqueResourceRepository[T]) GetByID(ctx context.Context, ID uint, preloads ...string) (*T, error) {
	tx, ok := getTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var resources T
	if err := tx.WithContext(ctx).Model(&resources).Where("id = ?", ID).First(&resources).Error; err != nil {
		return nil, err
	}

	return &resources, nil
}

func (r *userUniqueResourceRepository[T]) GetByIDs(ctx context.Context, ids []uint, preloads ...string) ([]T, error) {
	tx, ok := getTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var resources []T
	if err := tx.WithContext(ctx).Model(&resources).Where("id IN ?", ids).Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *userUniqueResourceRepository[T]) GetByUserID(ctx context.Context, userID uint, preloads ...string) ([]T, error) {
	tx, ok := getTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var resources []T
	if err := tx.WithContext(ctx).Model(&resources).Where("user_id = ?", userID).Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *userUniqueResourceRepository[T]) GetByCustom(ctx context.Context, costom string, preload ...string) ([]T, error) {
	tx, ok := getTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, p := range preload {
		tx = tx.Preload(p)
	}

	var resources []T
	if err := tx.Raw(costom).Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *userUniqueResourceRepository[T]) Where(ctx context.Context, param T, preloads ...string) ([]T, error) {
	tx, ok := getTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var resources []T
	if err := tx.Where(&param).Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *userUniqueResourceRepository[T]) Create(ctx context.Context, resource *T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	return tx.Create(resource).Error
}

func (r *userUniqueResourceRepository[T]) CreateOrUpdate(ctx context.Context, resource *T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	if !(*resource).IsEmpty() {
		var t T
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t, (*resource).GetID()).Error; err != nil {
			return err
		}
		return tx.Model(&t).Select((*resource).GetUpdateColumns()).Omit(clause.Associations).Updates(resource).Error
	}

	return tx.Create(resource).Error
}

func (r *userUniqueResourceRepository[T]) BulkCreate(ctx context.Context, resources []T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	return tx.Create(resources).Error
}

func (r *userUniqueResourceRepository[T]) SafeCreate(ctx context.Context, resources []T, update map[string]interface{}) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "resource_id"}},
		DoUpdates: clause.Assignments(update),
	}).Create(&resources).Error; err != nil {
		return err
	}

	return nil
}

func (r *userUniqueResourceRepository[T]) BulkUpdate(ctx context.Context, resources []T, param T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	ids := []uint{}
	for _, resource := range resources {
		ids = append(ids, (resource).GetID())
	}

	var t []T
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Find(&t).Error; err != nil {
		return err
	}

	return tx.Omit(clause.Associations).Where("id IN ?", ids).Updates(param).Error
}

func (r *userUniqueResourceRepository[T]) Updates(ctx context.Context, entities []T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	var ids []uint
	for _, entity := range entities {
		ids = append(ids, entity.GetID())
	}

	var t []T
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Find(&t).Error; err != nil {
		return err
	}

	for _, entity := range entities {
		if err := tx.Model(&entity).Select(entity.GetUpdateColumns()).
			Omit(clause.Associations).
			Updates(&entity).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userUniqueResourceRepository[T]) CustomUpdates(ctx context.Context, ids []uint, raw string) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	var t []T
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Find(&t).Error; err != nil {
		return err
	}

	return tx.Exec(raw).Error
}

func (r *userUniqueResourceRepository[T]) Delete(ctx context.Context, entity *T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	return tx.Unscoped().Delete(entity).Error
}

func (r *userUniqueResourceRepository[T]) BulkDelete(ctx context.Context, entities []T) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	return tx.Unscoped().Delete(entities).Error
}

func (r *userUniqueResourceRepository[T]) DeleteByUserID(ctx context.Context, userID uint) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	var entity T
	return tx.Unscoped().Where("user_id = ?", userID).Delete(&entity).Error
}
