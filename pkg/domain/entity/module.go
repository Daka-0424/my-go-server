package entity

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
