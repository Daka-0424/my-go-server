package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IUserPointSummary interface {
	Find(ctx context.Context, userID, platformNumber uint, paidKind int) (*entity.UserPointSummary, error)
	FirstOrCreateFreePointSummary(ctx context.Context, userID uint) (*entity.UserPointSummary, error)
	Update(ctx context.Context, pointSummary *entity.UserPointSummary) error
	BulkUpdate(ctx context.Context, points []entity.UserPointSummary) error
}
