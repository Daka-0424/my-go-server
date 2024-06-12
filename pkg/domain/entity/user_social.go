package entity

import "gorm.io/gorm"

type UserSocial struct {
	gorm.Model
	UserID uint
}
