package repository

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.User {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, Uuid, name string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil, repository.ErrTx
	}

	user := entity.User{Uuid: Uuid, Name: name}

	err := tx.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByUuid(ctx context.Context, uuid string, preloads ...string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	user := entity.User{}
	err := tx.First(&user, "uuid = ?", uuid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
