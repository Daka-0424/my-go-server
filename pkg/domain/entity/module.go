package entity

const DEFAULT_DB_ID = 0

func Entity() []any {
	return concatSlices([]any{
		// User
		&User{},
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
		UserResource(),
	)
}

func Seed() []any {
	return []any{
		&Item{},
		&PlatformProduct{},
	}
}

func UserResource() []any {
	return []any{
		&UserItem{},
	}
}

type ISeedType interface {
	SeedModule()
}

type IUserResourceType interface {
	UserResourceModule()
	GetID() uint
	IsEmpty() bool
}

func concatSlices(slices ...[]any) []any {
	var result []any
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
