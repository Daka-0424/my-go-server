package util

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("util",
		fx.Provide(
			NewLockFactoryUtil,
		))
}
