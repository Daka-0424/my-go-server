package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type User interface {
	CreateUser(ctx context.Context, Uuid, name, device, appVersion string, platform uint) (*entity.User, error)

	FindByUniqueUser(ctx context.Context, userId uint, uuid string, preloads ...string) (*entity.User, error)
	FindByUuid(ctx context.Context, uuid string, preloads ...string) (*entity.User, error)
	FindByUuids(ctx context.Context, uuids []string, preloads ...string) ([]*entity.User, error)
	FindByUserId(ctx context.Context, userId uint, preloads ...string) (*entity.User, error)
	FindByUserIds(ctx context.Context, userIds []uint, preloads ...string) ([]*entity.User, error)

	UpdateUser(ctx context.Context, user *entity.User) error
}
