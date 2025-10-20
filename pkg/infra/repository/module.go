package repository

import (
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"go.uber.org/fx"
)

func Modules() fx.Option {
	list := []interface{}{
		NewRedisCache,
		NewTransaction,
		NewUserRepository,
		NewUserSettingRepository,
		NewUserLoginStateRepository,
		NewUserPointSummaryRepository,
		NewEarnedPointRepository,
		NewUserSummaryRelationRepository,

		// Billing
		NewPaymentAppstoreTokenRepository,
		NewPaymentPlaystoreTokenRepository,

		// Admin
		NewAdminRepository,
	}

	// Seed
	list = append(list, seedRepositories()...)

	// User Resources
	list = append(list, userUserUniqueResourceRepositories()...)
	list = append(list, userResourcesRepositories()...)

	return fx.Module("repository",
		fx.Provide(
			list...,
		),
	)
}

func seedRepositories() []interface{} {
	return []interface{}{
		NewSeedRepository[entity.PlatformProduct],
	}
}

func userUserUniqueResourceRepositories() []interface{} {
	return []interface{}{
		NewUserUniqueResourceRepository[entity.UserItem],
	}
}

func userResourcesRepositories() []interface{} {
	return []interface{}{}
}
