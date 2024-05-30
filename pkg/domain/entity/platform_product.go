package entity

type PlatformProduct struct {
	SeedBase       `yaml:",inline"`
	Term           `yaml:",inline"`
	Name           string `yaml:"name" json:"name"`
	Description    string `yaml:"description" json:"description"`
	Price          uint   `yaml:"price" json:"price"`
	PlatformNumber uint   `yaml:"platformNumber" json:"platform_number"`
	PaidPoint      uint   `yaml:"paidPoint" json:"paid_point"`
	FreePoint      uint   `yaml:"freePoint" json:"free_point"`
	ProductId      string `yaml:"productId" json:"product_id"`
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
