package entity

type SeedBase struct {
	ID uint `yaml:"id" gorm:"primaryKey;autoIncrement:false"`
}
