package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IUserResource[T entity.IUserResourceType] interface {
	GetByID(ctx context.Context, ID uint, preloads ...string) (*T, error)
	GetByIDs(ctx context.Context, ISs []uint, preloads ...string) ([]T, error)
	GetByUserID(ctx context.Context, userID uint, preloads ...string) ([]T, error)
	Where(ctx context.Context, param T, preloads ...string) ([]T, error)

	CreateOrUpdate(ctx context.Context, entity *T) error
	BulkCreate(ctx context.Context, entities []T) error
	BulkUpdate(ctx context.Context, entities []T) error

	Delete(ctx context.Context, entity *T) error
	BulkDelete(ctx context.Context, entities []T) error
	DeleteByUserID(ctx context.Context, userID uint) error
}
