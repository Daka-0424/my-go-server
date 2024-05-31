package entity

import (
	"gorm.io/gorm"
)

type SpendPointRelation struct {
	gorm.Model
	UserID                   uint `json:"user_id"`
	UserPointSummaryID       uint `json:"user_point_summary_id"`
	EarnedPointID            uint `json:"earned_point_id"`
	SpendPointHistoryID      uint `json:"spend_point_history_id"`
	GemKind                  uint `json:"gem_kind"`
	PlatformNumberOnSpending uint `json:"platform_number_on_spending"`
	SpendPoint               uint `json:"spend_point"`
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
