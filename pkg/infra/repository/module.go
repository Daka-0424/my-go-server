package repository

import (
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Module("repository",
		fx.Provide(
			NewRedisCache,
			NewTransaction,
			NewUserRepository,

			// Seed
			NewSeedRepository[entity.VcPlatformProduct],
		),
	)
}
