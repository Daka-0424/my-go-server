package repository

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("repository",
		fx.Provide(
			NewRedisCache,
			NewTransaction,
			NewUserRepository,
		),
	)
}
