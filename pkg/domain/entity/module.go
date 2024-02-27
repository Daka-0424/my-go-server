package entity

const DEFAULT_DB_ID = 0

func Entity() []any {
	return []any{
		// User
		&User{},

		// Seed
		&PlatformProduct{},
	}
}

func Seed() []any {
	return []any{
		&PlatformProduct{},
	}
}

type SeedType interface {
	PlatformProduct
}
