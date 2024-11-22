package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userLoginStateRepository struct {
	db *gorm.DB
}

func NewUserLoginStateRepository(db *gorm.DB) repository.IUserLoginState {
	return &userLoginStateRepository{
		db: db,
	}
}

func (repo *userLoginStateRepository) CreateOrUpdate(ctx context.Context, state *entity.UserLoginState) error {
	tx, ok := GetTx(ctx)

	if !ok {
		return repository.ErrTx
	}

	if state.ID != 0 {
		t := entity.UserLoginState{Model: gorm.Model{ID: state.ID}}
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&t)
	}

	return tx.Omit(clause.Associations).Save(state).Error
}
