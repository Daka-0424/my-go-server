package entity

import (
	"time"

	"gorm.io/gorm"
)

type SpendPointHistory struct {
	gorm.Model
	UserID                   uint                `json:"user_id" gorm:"not null"`
	UserSummaryRelationID    uint                `json:"user_summary_relation_id" gorm:"not null"`
	UserSummaryRelation      UserSummaryRelation `json:"user_summary_relation"`
	PlatformNumberOnSpending uint                `json:"platform_number_on_spending" gorm:"not null"`
	SpendPoint               uint                `json:"spend_point" gorm:"default:0;not null"`
	SpendPaidPoint           uint                `json:"spend_paid_point" gorm:"default:0;not null"`
	SpendSalesAmount         uint                `json:"spend_sales_amount" gorm:"default:0;not null"`
	SpentAt                  time.Time           `json:"spent_at"`
	ItemCode                 string              `json:"item_code"`
}

func NewSpendPointHistory(userSummaryRelation *UserSummaryRelation, platformNumber, amount uint, itemCode string) *SpendPointHistory {
	return &SpendPointHistory{
		UserID:                   userSummaryRelation.UserID,
		UserSummaryRelationID:    userSummaryRelation.ID,
		UserSummaryRelation:      *userSummaryRelation,
		PlatformNumberOnSpending: platformNumber,
		SpendPoint:               amount,
		SpentAt:                  time.Now(),
		ItemCode:                 itemCode,
	}
}
