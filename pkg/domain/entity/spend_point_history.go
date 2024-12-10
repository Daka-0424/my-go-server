package entity

import (
	"time"

	"gorm.io/gorm"
)

type SpendPointHistory struct {
	gorm.Model
	UserID                   uint      `gorm:"user_id;not null"`
	UserSummaryRelationID    uint      `gorm:"user_summary_relation_id;not null"`
	PlatformNumberOnSpending uint      `gorm:"platform_number_on_spending;not null"`
	SpendPoint               uint      `gorm:"spend_point;default:0;not null"`
	SpendPaidPoint           uint      `gorm:"spend_paid_point;default:0;not null"`
	SpendSalesAmount         uint      `gorm:"spend_sales_amount;default:0;not null"`
	SpentAt                  time.Time `gorm:"spent_at"`
	ItemCode                 string    `gorm:"item_code"`
	UserSummaryRelation      UserSummaryRelation
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
