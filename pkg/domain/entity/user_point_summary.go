package entity

import (
	"gorm.io/gorm"
)

type UserPointSummary struct {
	gorm.Model
	UserID       uint `json:"user_id" gorm:"not null"`
	BalancePoint uint `json:"balance_point" gorm:"default:0;not null"`
	EarnPoint    uint `json:"earn_point" gorm:"default:0;not null"`
	PaidKind     int  `json:"paid_kind" gorm:"not null"`
	SpendPoint   uint `json:"spend_point" gorm:"default:0;not null"`
}

func NewUserPointSummary(
	userID uint,
	paidKind int,
) *UserPointSummary {
	return &UserPointSummary{
		UserID:   userID,
		PaidKind: paidKind,
	}
}

func (u *UserPointSummary) UpdateBalancePoint() {
	u.BalancePoint = u.EarnPoint - u.SpendPoint
}
