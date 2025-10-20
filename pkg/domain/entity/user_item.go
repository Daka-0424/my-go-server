package entity

type UserItem struct {
	UserUniqueResourceBase
	Resource Item `gorm:"foreignKey:ResourceID"`
	Quantity uint
}
