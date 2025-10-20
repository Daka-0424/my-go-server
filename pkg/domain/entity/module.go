package entity

const DEFAULT_DB_ID = 0

func Entity() []any {
	return concatSlices([]any{
		// User
		&User{},
		&UserSetting{},
		&UserLoginState{},
		&UserPointSummary{},
		&UserSummaryRelation{},

		// Admin
		&Admin{},

		// VC
		&EarnedPoint{},
		&ImitationPoint{},
		&SpendPointHistory{},
		&SpendPointRelation{},

		// Billing
		&PaymentAppstoreToken{},
		&PaymentPlaystoreToken{},
	},
		Seed(),
		UserUniqueResource(),
		UserResource(),
	)
}

func Seed() []any {
	return []any{
		&Item{},
		&PlatformProduct{},
	}
}

func UserUniqueResource() []any {
	return []any{
		// Item
		&UserItem{},
	}
}

func UserResource() []any {
	return []any{}
}

func concatSlices(slices ...[]any) []any {
	var result []any
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

type ISeedType interface {
	SeedModule()
}

type IUserResourceType interface {
	UserResourceModule()
	GetID() uint
	IsEmpty() bool
	GetUpdateColumns() []string
}

type UserUniqueResourceType interface {
	UserUniqueResourceModule()
	GetID() uint
	IsEmpty() bool
	GetUpdateColumns() []string
}
