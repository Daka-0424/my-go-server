package entity

import (
	"gorm.io/gorm"
)

type SpendPointRelation struct {
	gorm.Model
	UserID                   uint              `json:"user_id"`
	UserPointSummaryID       uint              `json:"user_point_summary_id"`
	UserPointSummary         UserPointSummary  `json:"user_point_summary"`
	EarnedPointID            uint              `json:"earned_point_id"`
	EarnedPoint              EarnedPoint       `json:"earned_point"`
	SpendPointHistoryID      uint              `json:"spend_point_history_id"`
	SpendPointHistory        SpendPointHistory `json:"spend_point_history"`
	GemKind                  uint              `json:"gem_kind"`
	PlatformNumberOnSpending uint              `json:"platform_number_on_spending"`
	SpendPoint               uint              `json:"spend_point"`
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
