package entity

import (
	"gorm.io/gorm"
)

type SpendPointRelation struct {
	gorm.Model
	UserID                   uint `gorm:"user_id"`
	UserPointSummaryID       uint `gorm:"user_point_summary_id"`
	EarnedPointID            uint `gorm:"earned_point_id"`
	SpendPointHistoryID      uint `gorm:"spend_point_history_id"`
	GemKind                  uint `gorm:"gem_kind"`
	PlatformNumberOnSpending uint `gorm:"platform_number_on_spending"`
	SpendPoint               uint `gorm:"spend_point"`
	UserPointSummary         UserPointSummary
	EarnedPoint              EarnedPoint
	SpendPointHistory        SpendPointHistory
}

func NewSpendPointRelation(userID, userPointSummaryID, earnedPointID, spendPointHistoryID,
	platformNumberOnSpending, spendPoint uint, gemKind uint) *SpendPointRelation {
	return &SpendPointRelation{
		UserID:                   userID,
		UserPointSummaryID:       userPointSummaryID,
		EarnedPointID:            earnedPointID,
		SpendPointHistoryID:      spendPointHistoryID,
		GemKind:                  gemKind,
		PlatformNumberOnSpending: platformNumberOnSpending,
		SpendPoint:               spendPoint,
	}
}
