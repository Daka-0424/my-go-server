package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	GemKindFree = iota
	GemKindPaid
)

type EarnedPoint struct {
	gorm.Model
	UserID                   uint             `json:"user_id"`
	UserPointSummaryID       uint             `json:"user_point_summary_id"`
	UserPointSummary         UserPointSummary `json:"user_point_summary"`
	PlatformProductID        uint             `json:"platform_product_id"`
	ImitationPointID         uint             `json:"imitation_point_id"`
	PlatformNumberOnSpending uint             `json:"platform_number_on_spending"`
	PointExceeded            bool             `json:"point_exceeded"`
	GemKind                  uint             `json:"gem_kind"`
	SpendPoint               uint             `json:"spend_point"`
	EarnedPoint              uint             `json:"earned_point"`
	BalancePoint             uint             `json:"balance_point"`
	EarnSource               string           `json:"earn_source"`
	SpentAt                  *time.Time       `json:"spent_at"`
	EarnedAt                 time.Time        `json:"earned_at"`
}

func NewCreateFreeEarnedPoint(vc *UserSummaryRelation, amount uint, vcPlatformProduct *PlatformProduct, earnSource string, imitationPoint *ImitationPoint) *EarnedPoint {
	// もしamountが0ならば、何もしない
	if amount <= 0 {
		return nil
	}

	earnedPoint := &EarnedPoint{
		UserID:                   vc.UserID,
		UserPointSummaryID:       vc.FreePointSummaryID,
		UserPointSummary:         vc.FreePointSummary,
		PlatformNumberOnSpending: uint(vc.PlatformNumber),
		GemKind:                  GemKindFree,
		EarnedPoint:              amount,
		BalancePoint:             amount,
		EarnSource:               earnSource,
		EarnedAt:                 time.Now(),
		SpentAt:                  func(d time.Time) *time.Time { return &d }(time.Now()),
	}

	// vcPlatformProductがnilでない場合に値をセット
	if vcPlatformProduct != nil {
		earnedPoint.PlatformProductID = vcPlatformProduct.ID
	}

	// imitationPointがnilでない場合に値をセット
	if imitationPoint != nil {
		earnedPoint.ImitationPointID = imitationPoint.ID
	}

	return earnedPoint
}

func NewCreatePaidEarnedPoint(vc *UserSummaryRelation, amount uint, vcPlatformProduct *PlatformProduct, earnSource string, imitationPoint *ImitationPoint, earnedAt time.Time) *EarnedPoint {
	// productionでは、レシートも補填配信操作も経由しない有償発行は許可しない

	earnPoint := vcPlatformProduct.PaidPoint
	if amount > 0 {
		earnPoint = amount
	}

	earnedPoint := &EarnedPoint{
		UserID:                   vc.UserID,
		UserPointSummaryID:       vc.PaidPointSummaryID,
		UserPointSummary:         vc.PaidPointSummary,
		PlatformProductID:        vcPlatformProduct.ID,
		PlatformNumberOnSpending: uint(vc.PlatformNumber),
		GemKind:                  GemKindPaid,
		EarnedPoint:              earnPoint,
		BalancePoint:             earnPoint,
		EarnSource:               earnSource,
		EarnedAt:                 earnedAt,
	}

	// imitationPointがnilでない場合に値をセット
	if imitationPoint != nil {
		earnedPoint.ImitationPointID = imitationPoint.ID
	}

	return earnedPoint
}

func (e *EarnedPoint) UnitCost() float64 {
	if e.GemKind == GemKindFree {
		return 0.0
	}
	if e.PlatformProductID == DEFAULT_DB_ID {
		return 0.0
	}
	return 0.0
	//return e.PlatformProduct.UnitCost()
}

func (e *EarnedPoint) UpdateBalancePoint() {
	e.BalancePoint = e.EarnedPoint - e.SpendPoint
	if e.BalancePoint == 0 {
		e.PointExceeded = true
	}
}
