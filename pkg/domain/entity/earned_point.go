package entity

import "gorm.io/gorm"

type EarnedPoint struct {
	gorm.Model
	UserID                  uint
	User                    User
	UserPointSummaryID      uint
	UserPointSummary        UserPointSummary
	PaidKind                uint
	PlatformNumberOnEarnign uint
}
