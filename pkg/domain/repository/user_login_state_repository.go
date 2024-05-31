package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IUserLoginState interface {
	CreateOrUpdate(ctx context.Context, state *entity.UserLoginState) error
}
