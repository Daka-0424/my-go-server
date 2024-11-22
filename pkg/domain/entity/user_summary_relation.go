package entity

import (
	"gorm.io/gorm"
)

type UserSummaryRelation struct {
	gorm.Model
	UserID             uint `gorm:"user_id" gorm:"not null <-:create"`
	PlatformNumber     uint `gorm:"platform_number" gorm:"not null"`
	FreePointSummaryID uint `gorm:"free_point_summary_id" gorm:"not null"`
	PaidPointSummaryID uint `gorm:"paid_point_summary_id" gorm:"not null"`
	FreePointSummary   UserPointSummary
	PaidPointSummary   UserPointSummary
}

func NewUserSummaryRelation(
	userID uint,
	platformNumber uint,
) *UserSummaryRelation {
	return &UserSummaryRelation{
		UserID:         userID,
		PlatformNumber: platformNumber,
	}
}

func (u *UserSummaryRelation) PointSummaries() []UserPointSummary {
	return []UserPointSummary{u.FreePointSummary, u.PaidPointSummary}
}

func (u *UserSummaryRelation) BalancePoint() uint {
	return u.FreePointSummary.BalancePoint + u.PaidPointSummary.BalancePoint
}

func (u *UserSummaryRelation) BalanceFreePoint() uint {
	return u.FreePointSummary.BalancePoint
}

func (u *UserSummaryRelation) BalancePaidPoint() uint {
	return u.PaidPointSummary.BalancePoint
}

func (u *UserSummaryRelation) EarnedPoint() uint {
	return u.FreePointSummary.EarnPoint + u.PaidPointSummary.EarnPoint
}

func (u *UserSummaryRelation) EarnedFreePoint() uint {
	return u.FreePointSummary.EarnPoint
}

func (u *UserSummaryRelation) EarnedPaidPoint() uint {
	return u.PaidPointSummary.EarnPoint
}

func (u *UserSummaryRelation) SpentPoint() uint {
	return u.FreePointSummary.SpendPoint + u.PaidPointSummary.SpendPoint
}

func (u *UserSummaryRelation) SpentFreePoint() uint {
	return u.FreePointSummary.SpendPoint
}

func (u *UserSummaryRelation) SpentPaidPoint() uint {
	return u.PaidPointSummary.SpendPoint
}
