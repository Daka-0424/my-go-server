package entity

type PlatformProduct struct {
	SeedBase       `yaml:",inline"`
	Term           `yaml:",inline"`
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	Price          uint   `yaml:"price"`
	PlatformNumber uint   `yaml:"platformNumber"`
	PaidPoint      uint   `yaml:"paidPoint"`
	FreePoint      uint   `yaml:"freePoint"`
	ProductId      string `yaml:"productId"`
}
