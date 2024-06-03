package entity

import "gorm.io/gorm"

type UserItem struct {
	gorm.Model
	UserRewardContent
	Item Item `gorm:"foreignKey:ContentID"`
}
