package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

const (
	PreloadResourceGrowthFundLineup = "GrowthFundLineup"
)

type IUserUniqueResource[T entity.UserUniqueResourceType] interface {
	GetByID(ctx context.Context, ID uint, preloads ...string) (*T, error)
	GetByIDs(ctx context.Context, ISs []uint, preloads ...string) ([]T, error)
	GetByUserID(ctx context.Context, userID uint, preloads ...string) ([]T, error)
	GetByCustom(ctx context.Context, costom string, preload ...string) ([]T, error)
	Where(ctx context.Context, param T, preloads ...string) ([]T, error)
	Create(ctx context.Context, resource *T) error
	CreateOrUpdate(ctx context.Context, resource *T) error
	BulkCreate(ctx context.Context, resources []T) error
	SafeCreate(ctx context.Context, resources []T, update map[string]interface{}) error
	BulkUpdate(ctx context.Context, resources []T, param T) error
	Updates(ctx context.Context, entities []T) error
	CustomUpdates(ctx context.Context, ids []uint, raw string) error
	Delete(ctx context.Context, entity *T) error
	BulkDelete(ctx context.Context, entities []T) error
	DeleteByUserID(ctx context.Context, userID uint) error
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
