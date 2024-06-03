package entity

type UserRewardContent struct {
	RewardContent
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`
}
