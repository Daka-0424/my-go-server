package entity

import "gorm.io/gorm"

type UserUniqueResourceBase struct {
	gorm.Model
	UserID     uint `gorm:"<-:create;not null;uniqueIndex:idx_user_id_resource_id,priority:1"`
	ResourceID uint `gorm:"<-:create;not null;uniqueIndex:idx_user_id_resource_id,priority:2"`

	gormAuxiliary
}

func (r UserUniqueResourceBase) UserUniqueResourceModule() {}

func (r UserUniqueResourceBase) GetID() uint {
	return r.ID
}

func (r UserUniqueResourceBase) IsEmpty() bool {
	return r.ID == 0
}

func (r UserUniqueResourceBase) GetUpdateColumns() []string {
	return r.gormAuxiliary.GetUpdateColumns()
}
