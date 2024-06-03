package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IUserResource[T entity.IUserResourceType] interface {
	GetByID(ctx context.Context, ID uint) (*T, error)
	GetByIDs(ctx context.Context, IDs []uint) ([]T, error)
	GetByUserID(ctx context.Context, userID uint) ([]T, error)
	CreateOrUpdate(ctx context.Context, resource *T) error
	BulkCreate(ctx context.Context, resources []T) error
	BulkUpdate(ctx context.Context, resources []T) error
}
