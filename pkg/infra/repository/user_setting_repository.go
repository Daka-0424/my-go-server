package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) repository.IUserSetting {
	return &userSettingRepository{
		db: db,
	}
}

func (r *userSettingRepository) CreateOrUpdate(ctx context.Context, setting *entity.UserSetting) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	if setting.ID != 0 {
		t := entity.UserSetting{Model: gorm.Model{ID: setting.ID}}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t).Error; err != nil {
			return err
		}
		return tx.Model(setting).Select(setting.GetUpdateColumns()).Omit(clause.Associations).Updates(setting).Error
	}

	return tx.Create(setting).Error
}

func (r *userSettingRepository) Delete(ctx context.Context, setting *entity.UserSetting) error {
	tx, ok := getTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	return tx.Unscoped().Delete(setting).Error
}
