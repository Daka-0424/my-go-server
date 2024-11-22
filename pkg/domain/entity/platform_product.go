package entity

type PlatformProduct struct {
	SeedBase       `yaml:",inline"`
	Term           `yaml:",inline"`
	Name           string `yaml:"name" gorm:"name"`
	Description    string `yaml:"description" gorm:"description"`
	Price          uint   `yaml:"price" gorm:"price"`
	PlatformNumber uint   `yaml:"platformNumber" gorm:"platform_number"`
	PaidPoint      uint   `yaml:"paidPoint" gorm:"paid_point"`
	FreePoint      uint   `yaml:"freePoint" gorm:"free_point"`
	ProductId      string `yaml:"productId" gorm:"product_id"`
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
