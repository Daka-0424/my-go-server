package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type User interface {
	CreateUser(ctx context.Context, Uuid, name string) (*entity.User, error)

	FindByUuid(ctx context.Context, uuid string, preloads ...string) (*entity.User, error)
}
