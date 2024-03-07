package service

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("service",
		fx.Provide(
			NewUserService,
			NewVcService,
		),
	)
}
