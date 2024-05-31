package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IUserSummaryRelation interface {
	FindByUserID(ctx context.Context, userID uint) (*entity.UserSummaryRelation, error)
	CreateOrUpdate(ctx context.Context, entity *entity.UserSummaryRelation) error
}
