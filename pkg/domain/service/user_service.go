package service

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

const USER_DEFAULT_NAME = "NewUser"

type User interface {
	CreateUser(ctx context.Context, uuid, device, appVersion string, platform uint) (*entity.User, error)
}

type userService struct {
	userRepository repository.User
}

func NewUserService(ur repository.User) User {
	return &userService{
		userRepository: ur,
	}
}

func (u *userService) CreateUser(ctx context.Context, uuid, device, appVersion string, platform uint) (*entity.User, error) {
	user, err := u.userRepository.FindByUuid(ctx, uuid)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("user already exist")
	}

	return u.userRepository.CreateUser(ctx, uuid, USER_DEFAULT_NAME, device, appVersion, platform)
}
