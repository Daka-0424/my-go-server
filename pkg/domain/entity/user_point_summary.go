package entity

import (
	"gorm.io/gorm"
)

type UserPointSummary struct {
	gorm.Model
	UserID       uint `gorm:"user_id;not null"`
	BalancePoint uint `gorm:"balance_point;default:0;not null"`
	EarnPoint    uint `gorm:"earn_point;default:0;not null"`
	PaidKind     int  `gorm:"paid_kind;not null"`
	SpendPoint   uint `gorm:"spend_point;default:0;not null"`

	gormAuxiliary
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

func (u *UserPointSummary) AddEarnPoint(amount uint) {
	u.EarnPoint += amount
	u.updateColumn("earn_point")
}

func (u *UserPointSummary) UpdateBalancePoint() {
	u.BalancePoint = u.EarnPoint - u.SpendPoint
	u.updateColumn("balance_point")
}
