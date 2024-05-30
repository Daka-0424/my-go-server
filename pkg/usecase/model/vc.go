package model

import "github.com/Daka-0424/my-go-server/pkg/domain/entity"

type Vc struct {
	UserID                uint `json:"user_id"`
	UserSummaryRelationID uint `json:"user_summary_relation_id"`
	BalancePoint          uint `json:"balance_point"`
	BalanceFreePoint      uint `json:"balance_free_point"`
	BalancePaidPoint      uint `json:"balance_paid_point"`
	EarnedPoint           uint `json:"earned_point"`
	EarnedFreePoint       uint `json:"earned_free_point"`
	EarnedPaidPoint       uint `json:"earned_paid_point"`
	SpentPoint            uint `json:"spent_point"`
	SpentFreePoint        uint `json:"spent_free_point"`
	SpentPaidPoint        uint `json:"spent_paid_point"`
}
type ReceiptResult struct {
	Vc                *Vc
	PlatformProductID uint `json:"platform_product_id"`
}

func NewVc(userVc *entity.UserSummaryRelation) *Vc {
	return &Vc{
		UserID:                userVc.UserID,
		UserSummaryRelationID: userVc.ID,
		BalancePoint:          userVc.BalancePoint(),
		BalanceFreePoint:      userVc.BalanceFreePoint(),
		BalancePaidPoint:      userVc.BalancePaidPoint(),
		EarnedPoint:           userVc.EarnedPoint(),
		EarnedFreePoint:       userVc.EarnedFreePoint(),
		EarnedPaidPoint:       userVc.EarnedPaidPoint(),
		SpentPoint:            userVc.SpentPoint(),
		SpentFreePoint:        userVc.SpentFreePoint(),
		SpentPaidPoint:        userVc.SpentPaidPoint(),
	}
}

func NewReceiptResult(user *entity.UserSummaryRelation, platformProduct *entity.PlatformProduct) *ReceiptResult {
	return &ReceiptResult{
		Vc:                NewVc(user),
		PlatformProductID: platformProduct.ID,
	}
}
