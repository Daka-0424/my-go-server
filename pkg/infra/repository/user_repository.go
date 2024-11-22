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

func NewUserRepository(db *gorm.DB) repository.IUser {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) ExistsUser(ctx context.Context, uuid string) (bool, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		return false, repository.ErrTx
	}

	var count int64
	if err := tx.Model(&entity.User{}).Where("uuid = ?", uuid).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, Uuid, name, device, clientVersion string, platform uint) (*entity.User, error) {
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

func (repo *userRepository) CreateUserSummaryRelation(ctx context.Context, vc *entity.UserSummaryRelation) error {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	return tx.Create(&vc).Error
}

func (repo *userRepository) FindByUniqueUser(ctx context.Context, userId uint, uuid string, preloads ...string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = repo.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	user := entity.User{}
	err := tx.First(&user, "id = ? AND uuid = ?", userId, uuid).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) FindByUuid(ctx context.Context, uuid string, preloads ...string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = repo.db
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

func (repo *userRepository) FindByUuids(ctx context.Context, uuids []string, preloads ...string) ([]*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = repo.db
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

func (repo *userRepository) FindByUserId(ctx context.Context, userId uint, preloads ...string) (*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = repo.db
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

func (repo *userRepository) FindByUserIds(ctx context.Context, userIds []uint, preloads ...string) ([]*entity.User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		tx = repo.db
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

func (repo *userRepository) FindUserWithVc(ctx context.Context, userID uint) (*entity.User, *entity.UserSummaryRelation, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var user entity.User
	var vc entity.UserSummaryRelation

	if err := tx.First(&user, "id = ?", userID).Error; err != nil {
		return &entity.User{}, &entity.UserSummaryRelation{}, err
	}

	if err := tx.Where("user_id = ?", userID).
		Preload("PaidPointSummary").
		Preload("FreePointSummary").
		First(&vc).Error; err != nil {
		return &user, &entity.UserSummaryRelation{}, err
	}

	return &user, &vc, nil
}

func (repo *userRepository) FindUserPointSummary(ctx context.Context, userID uint, platformNumber uint, paidKind int) (*entity.UserPointSummary, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var vc entity.UserPointSummary

	if err := tx.Where("user_id = ? AND platform_number = ? AND paid_kind = ?", userID, platformNumber, paidKind).
		First(&vc).Error; err != nil {
		return nil, err
	}

	return &vc, nil
}

func (repo *userRepository) FindOtherPlatformVc(ctx context.Context, userID uint, platformNumber uint) (*entity.UserSummaryRelation, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var vc entity.UserSummaryRelation

	if err := tx.Where("user_id = ? AND platform_number != ?", userID, platformNumber).
		First(&vc).Error; err != nil {
		return nil, err
	}

	return &vc, nil
}

func (repo *userRepository) FirstOrCreateFreePointSummary(ctx context.Context, userID uint, paidKind int) (*entity.UserPointSummary, error) {
	tx, ok := GetTx(ctx)

	if !ok {
		tx = repo.db
	}

	var vc entity.UserPointSummary
	if err := tx.FirstOrCreate(&vc, entity.UserPointSummary{UserID: userID, PaidKind: 0}).Error; err != nil {
		return nil, err
	}

	return &vc, nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return repository.ErrTx
	}

	t := entity.User{Model: gorm.Model{ID: user.ID}}
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t)

	return tx.Model(&entity.User{}).Where("id = ?", user.ID).Select("*").Omit(clause.Associations).Updates(user).Error
}
