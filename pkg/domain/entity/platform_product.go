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

func (p *PlatformProduct) UnitCost() float64 {
	if p.PrepaidPaymentMethod() {
		return float64(p.Price) / float64(p.PaidPoint)
	} else {
		return float64(p.Price)
	}
}

func (p *PlatformProduct) PrepaidPaymentMethod() bool {
	return p.PaidPoint > 0
}
