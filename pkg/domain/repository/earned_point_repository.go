package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IEarnedPoint interface {
	CreateOrUpdate(ctx context.Context, earnedPoint *entity.EarnedPoint) error
	BulkCreate(ctx context.Context, earnedPoints []entity.EarnedPoint) error
	GetAll(ctx context.Context, offset int, limit int) ([]entity.EarnedPoint, error)
	GetWhere(ctx context.Context, param entity.EarnedPoint, offset int, limit int) ([]entity.EarnedPoint, error)
	FindByPointSummaryIDs(ctx context.Context, pointSummaryIDs ...uint) ([]entity.EarnedPoint, error)
	CountAll(ctx context.Context) (int64, error)
	CountWhere(ctx context.Context, param entity.EarnedPoint) (int64, error)
}
