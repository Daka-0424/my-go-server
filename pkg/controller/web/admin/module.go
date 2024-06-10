package admin

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("admin",
		fx.Provide(
			NewAdminController,
		),
	)
}
