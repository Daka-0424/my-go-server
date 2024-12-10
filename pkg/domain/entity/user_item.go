package entity

type UserItem struct {
	UserResourceBase
	Resource Item `gorm:"foreignKey:ResourceID"`
	Quantity uint `gorm:"quantity;not null;index:idx_user_id_resource_id,priority:3"`
}
