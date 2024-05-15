package appstore

import (
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Module("appstore",
		fx.Provide(
			NewAppStoreFactory,
		),
	)
}
