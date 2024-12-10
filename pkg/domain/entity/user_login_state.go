package entity

import (
	"time"

	"github.com/Songmu/flextime"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

type UserLoginState struct {
	gorm.Model
	UserID          uint       `gorm:"user_id;not null"`
	TotalLogin      uint       `gorm:"total_login;not null"`
	Duration        *uint      `gorm:"duration;not null"`
	AccessedAt      *time.Time `gorm:"accessed_at"`
	DurationStartAt *time.Time `gorm:"duration_start_at"`
	LastLoginAt     *time.Time `gorm:"last_login_at"`
}

func NewUserLoginState(userID uint) *UserLoginState {
	return &UserLoginState{
		UserID:     userID,
		Duration:   func(d uint) *uint { return &d }(0),
		AccessedAt: func(t time.Time) *time.Time { return &t }(flextime.Now()),
	}
}

func (login *UserLoginState) Login(user *User) bool {
	today := flextime.Now()
	login.AccessedAt = &today

	if login.HasLoggedInToday(user) {
		return false
	}

	yesterday := today.AddDate(0, 0, -1)

	beginingOfDay := now.With(yesterday).BeginningOfDay()
	endOfDay := now.With(yesterday).EndOfDay()

	if login.LastLoginAt != nil && !beginingOfDay.After(*login.LastLoginAt) && (*login.LastLoginAt).Before(endOfDay) {
		*login.Duration++
	} else {
		*login.Duration = 1
		login.DurationStartAt = &today
	}

	login.TotalLogin++
	login.LastLoginAt = &today

	return true
}

func (login *UserLoginState) HasLoggedInToday(user *User) bool {
	time := flextime.Now()
	if login.LastLoginAt == nil {
		return false
	}

	begin := now.With(*login.LastLoginAt).BeginningOfDay()
	end := now.With(*login.LastLoginAt).EndOfDay()

	return !begin.After(time) && time.Before(end)
}
