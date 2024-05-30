package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

const (
	USER_DEFAULT_NAME = "NewUser"
	CodeKey           = 81
)

type IUser interface {
	Register(ctx context.Context, uuid, device, clientVersion string, platformNumber uint) (*entity.User, error)
}

type userService struct {
	userRepository           repository.IUser
	userLoginStateRepository repository.IUserLoginState
}

func NewUserService(
	ur repository.IUser,
	ulsr repository.IUserLoginState,
) IUser {
	return &userService{
		userRepository:           ur,
		userLoginStateRepository: ulsr,
	}
}

func (service *userService) Register(ctx context.Context, uuid, device, clientVersion string, platformNumber uint) (*entity.User, error) {
	user, err := service.userRepository.CreateUser(ctx, uuid, USER_DEFAULT_NAME, clientVersion, device, platformNumber)
	if err != nil {
		return nil, err
	}

	// DisplayCodeを作成
	user.DisplayCode = service.createDisplayCode(user)
	if err := service.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// その後、userを返す
	return user, nil
}

func (service *userService) createDisplayCode(user *entity.User) string {
	first := rune((user.CreatedAt.Year() - 2020 + 45) % 256)
	second := rune((int(user.CreatedAt.Month()) + 67) % 256)
	code := hash(user.ID)

	return fmt.Sprintf("%c%c%02d", first, second, code)
}

func hash(userID uint) int {
	multiplier := rand.Intn(CodeKey)
	h := int(userID)
	h = h*multiplier + CodeKey
	return h
}
