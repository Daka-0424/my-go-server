package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type User interface {
	ExistsUser(ctx context.Context, uuid string) (bool, error)

	CreateUser(ctx context.Context, Uuid, name, device, clientVersion string, platform uint) (*entity.User, error)
	CreateUserParams(ctx context.Context, userID uint) error
	CreateUserSummaryRelation(ctx context.Context, vc *entity.UserSummaryRelation) error

	FindByUniqueUser(ctx context.Context, userId uint, uuid string, preloads ...string) (*entity.User, error)
	FindByUuid(ctx context.Context, uuid string, preloads ...string) (*entity.User, error)
	FindByUuids(ctx context.Context, uuids []string, preloads ...string) ([]*entity.User, error)
	FindByUserId(ctx context.Context, userId uint, preloads ...string) (*entity.User, error)
	FindByUserIds(ctx context.Context, userIds []uint, preloads ...string) ([]*entity.User, error)
	FindUserWithVc(ctx context.Context, userID uint) (*entity.User, *entity.UserSummaryRelation, error)

	FindUserPointSummary(ctx context.Context, userID uint, platformNumber uint, paidKind int) (*entity.UserPointSummary, error)
	FindOtherPlatformVc(ctx context.Context, userID uint, platformNumber uint) (*entity.UserSummaryRelation, error)
	FirstOrCreateFreePointSummary(ctx context.Context, userID uint, paidKind int) (*entity.UserPointSummary, error)

	UpdateUser(ctx context.Context, user *entity.User) error
}
