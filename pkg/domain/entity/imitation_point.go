package entity

import "gorm.io/gorm"

type ImitationPoint struct {
	gorm.Model
	UserID                uint                `json:"user_id" gorm:"not null"`
	UserSummaryRelationID uint                `json:"user_summary_relation_id" gorm:"not null"`
	UserSummaryRelation   UserSummaryRelation `json:"user_summary_relation"`
	PlatformProductID     uint                `json:"platform_product_id"`
	FreeEarnedPointID     uint                `json:"free_earned_point_id"`
	PaidEarnedPointID     uint                `json:"paid_earned_point_id"`
	AdminUserID           uint                `json:"admin_user_id"`
	AdminUser             Admin               `json:"admin_user"`
	WithFreePoint         uint                `json:"with_free_point"`
	//WithReportEntry       uint                `json:"with_report_entry"`
	Comment string `json:"comment"`
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
