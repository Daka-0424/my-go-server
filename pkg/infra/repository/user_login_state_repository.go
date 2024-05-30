package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userLoginStateRepository struct {
	db     *gorm.DB
	fields []string
}

func NewUserLoginStateRepository(db *gorm.DB) repository.IUserLoginState {
	return &userLoginStateRepository{
		db:     db,
		fields: entity.GetEntityFields(entity.UserLoginState{}),
	}
}

func (repo *userLoginStateRepository) CreateOrUpdate(ctx context.Context, state *entity.UserLoginState) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},         // 一意性を保証するカラム
		DoUpdates: clause.AssignmentColumns(repo.fields), // 更新するフィールドを指定
	}).Create(state).Error
}
