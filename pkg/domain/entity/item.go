package entity

type Item struct {
	SeedBase `yaml:",inline"`
	Name     string `yaml:"name"`
}
