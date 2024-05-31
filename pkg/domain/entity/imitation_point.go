package entity

import (
	"gorm.io/gorm"
)

type ImitationPoint struct {
	gorm.Model
	UserID                uint   `json:"user_id" gorm:"not null"`
	UserSummaryRelationID uint   `json:"user_summary_relation_id" gorm:"not null"`
	PlatformProductID     uint   `json:"platform_product_id" gorm:"not null"`
	FreeEarnedPointID     uint   `json:"free_earned_point_id"`
	PaidEarnedPointID     uint   `json:"paid_earned_point_id"`
	AdminUserID           uint   `json:"admin_user_id" gorm:"not null"`
	WithFreePoint         uint   `json:"with_free_point"`
	Comment               string `json:"comment"`
	UserSummaryRelation   UserSummaryRelation
	PlatformProduct       PlatformProduct
	AdminUser             Admin
}

func NewImitationPoint(userID, userSummaryRelationID, platformProductID, adminID, withFreePoint uint, comment string) *ImitationPoint {
	return &ImitationPoint{
		UserID:                userID,
		UserSummaryRelationID: userSummaryRelationID,
		PlatformProductID:     platformProductID,
		AdminUserID:           adminID,
		WithFreePoint:         withFreePoint,
		Comment:               comment,
	}
}
