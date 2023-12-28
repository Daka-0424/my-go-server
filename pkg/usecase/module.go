package usecase

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("usecase",
		fx.Provide(
			NewRegistrationUsecase,
			NewSessionUsecase,
		),
	)
}
