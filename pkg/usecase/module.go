package usecase

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Module("usecase",
		fx.Provide(
			NewUserUsecase,
			NewSessionUsecase,
			NewPlatformProductUsecase,

			// Billing
			NewAppstoreUsecase,
			NewPlaystoreUsecase,

			// Admin
			NewAdminUsecase,
		),
	)
}
