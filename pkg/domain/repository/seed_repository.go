package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type ISeed[T entity.ISeedType] interface {
	GetByID(ctx context.Context, ID uint) (*T, error)
	GetByIDs(ctx context.Context, IDs []uint) ([]T, error)
	GetAll(ctx context.Context, preloads ...string) ([]T, error)
	Get(ctx context.Context, param T) (*T, error)
	Where(ctx context.Context, param T) ([]T, error)
	GetByMap(ctx context.Context, param map[string]interface{}) (*T, error)
}
