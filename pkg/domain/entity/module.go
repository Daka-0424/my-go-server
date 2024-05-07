package entity

const DEFAULT_DB_ID = 0

func Entity() []any {
	return []any{
		// User
		&User{},
	}
}

func ReadEntity() []any {
	return []any{
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
