package controller

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("controller",
		fx.Provide(
			NewUserController,
			NewSessionController,
			NewPlatformProductController,
		),
	)
}
