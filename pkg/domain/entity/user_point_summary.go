package entity

import "gorm.io/gorm"

const (
	FreePoint = iota // 無償
	Paid             // 有償
)

type UserPointSummary struct {
	gorm.Model
	UserID uint
	User User
	PointKind uint
}