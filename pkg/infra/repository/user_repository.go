package repository

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.User {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, Uuid, name, device, clientVersion string, platform uint) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil, repository.ErrTx
	}

	user := entity.User{
		UUID:           Uuid,
		Name:           name,
		ClientVersion:  clientVersion,
		Device:         device,
		PlatformNumber: platform,
	}

	err := tx.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUserParams(ctx context.Context, userID uint) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	loginState := entity.NewUserLoginState(userID)
	err := tx.Create(&loginState).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByUniqueUser(ctx context.Context, userId uint, uuid string, preloads ...string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	user := entity.User{}
	err := tx.First(&user, "id = ? AND uuid = ?", userId, uuid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
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

func (r *userRepository) FindByUuids(ctx context.Context, uuids []string, preloads ...string) ([]*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var users []*entity.User
	err := tx.Where("uuid IN ?", uuids).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) FindByUserId(ctx context.Context, userId uint, preloads ...string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	user := entity.User{}
	err := tx.First(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByUserIds(ctx context.Context, userIds []uint, preloads ...string) ([]*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var users []*entity.User
	err := tx.Where("id IN ?", userIds).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	t := entity.User{Model: gorm.Model{ID: user.ID}}
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t)

	return tx.Model(user).Select("*").Omit(clause.Associations).Updates(user).Error
}
