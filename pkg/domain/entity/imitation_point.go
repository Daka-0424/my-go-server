package entity

import (
	"gorm.io/gorm"
)

type ImitationPoint struct {
	gorm.Model
	UserID                uint   `gorm:"user_id" gorm:"not null"`
	UserSummaryRelationID uint   `gorm:"user_summary_relation_id" gorm:"not null"`
	PlatformProductID     uint   `gorm:"platform_product_id" gorm:"not null"`
	FreeEarnedPointID     uint   `gorm:"free_earned_point_id"`
	PaidEarnedPointID     uint   `gorm:"paid_earned_point_id"`
	AdminUserID           uint   `gorm:"admin_user_id" gorm:"not null"`
	WithFreePoint         uint   `gorm:"with_free_point"`
	Comment               string `gorm:"comment"`
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
