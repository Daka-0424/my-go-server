package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
	fields []string
}

func NewAdminRepository(db *gorm.DB) repository.IAdmin {
	return &adminRepository{
		db: db,
		fields: entity.GetEntityFields(entity.Admin{}),
	}
}

func (repo *adminRepository) Exsists(ctx context.Context, email string) bool {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = repo.db
	}

	admin := entity.Admin{}
	err := tx.First(&admin, "email = ?", email).Error

	return err == nil
}

func (repo *adminRepository) Register(ctx context.Context, email, pass string, roleType entity.AdminRoleType) (*entity.Admin, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil, repository.ErrTx
	}

	admin := entity.NewAdmin(email, pass, roleType)
	if err := tx.Create(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

func (repo *adminRepository) Find(ctx context.Context, param entity.Admin) ([]entity.Admin, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil, repository.ErrTx
	}

	admins := []entity.Admin{}
	if err := tx.Where(param).Find(&admins).Error; err != nil {
		return nil, err
	}

	return admins, nil
}

func (repo *adminRepository) Update(ctx context.Context, admin *entity.Admin) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	if err := tx.Save(admin).Error; err != nil {
		return err
	}

	return nil
}

func (repo *adminRepository) Delete(ctx context.Context, admin *entity.Admin) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	if err := tx.Unscoped().Delete(admin).Error; err != nil {
		return err
	}

	return nil
}

func (repo *adminRepository) CountAll(ctx context.Context) (int64, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		return 0, repository.ErrTx
	}

	var count int64
	if err := tx.Model(&entity.Admin{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
