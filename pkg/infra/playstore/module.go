package playstore

import (
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Module("playstore",
		fx.Provide(
			NewPlaystoreFactory,
		),
	)
}
