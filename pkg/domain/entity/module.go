package entity

func Entity() []any {
	return []any{
		// User
		&User{},

		// Seed
		&VcPlatformProduct{},
	}
}

func Seed() []any {
	return []any{
		&VcPlatformProduct{},
	}
}

type SeedType interface {
	VcPlatformProduct
}
