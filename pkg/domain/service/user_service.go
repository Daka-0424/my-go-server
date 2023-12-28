package service

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type User interface {
	CreateUser(ctx context.Context, uuid, name string) (*entity.User, error)
}

type userService struct {
	userRepository repository.User
}

func NewUserService(ur repository.User) User {
	return &userService{
		userRepository: ur,
	}
}

func (u *userService) CreateUser(ctx context.Context, uuid, name string) (*entity.User, error) {
	user, err := u.userRepository.FindByUuid(ctx, uuid)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("user already exist")
	}

	return u.userRepository.CreateUser(ctx, uuid, name)
}
