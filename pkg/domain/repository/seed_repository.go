package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type ISeed[T entity.ISeedType] interface {
	GetAll(ctx context.Context, preloads ...string) ([]T, error)
	Where(ctx context.Context, param T) ([]T, error)
}
