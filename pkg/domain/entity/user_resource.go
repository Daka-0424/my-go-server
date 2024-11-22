package entity

import "gorm.io/gorm"

type UserResourceBase struct {
	gorm.Model
	UserID     uint `gorm:"not null;index:idx_user_id_resource_id,priority:1"`
	User       User `gorm:"foreignKey:UserID"`
	ResourceID uint `gorm:"not null;index:idx_user_id_resource_id,priority:2"`
}

func (r UserResourceBase) UserResourceModule() {}

func (r UserResourceBase) GetID() uint {
	return r.ID
}

func (r UserResourceBase) IsEmpty() bool {
	return r.ID == 0
}
