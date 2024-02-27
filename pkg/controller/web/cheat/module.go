package cheat

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("cheat",
		fx.Provide(
			NewCheatRootController,
		),
	)
}
