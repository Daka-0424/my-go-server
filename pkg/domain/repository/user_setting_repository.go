package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IUserSetting interface {
	CreateOrUpdate(ctx context.Context, setting *entity.UserSetting) error
	Delete(ctx context.Context, setting *entity.UserSetting) error
}
