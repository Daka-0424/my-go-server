package entity

import (
	"time"

	"github.com/Songmu/flextime"
)

type Term struct {
	StartAt     DateTime `yaml:"startAt" gorm:"start_at"`
	EndAt       DateTime `yaml:"endAt" gorm:"end_at"`
	TestStartAt DateTime `yaml:"testStartAt" gorm:"test_start_at"`
	TestEndAt   DateTime `yaml:"testEndAt" gorm:"test_end_at"`
}

func NewTerm(startAt, endAt, testStartAt, testEndAt time.Time) *Term {
	return &Term{
		StartAt:     DateTime{startAt},
		EndAt:       DateTime{endAt},
		TestStartAt: DateTime{testStartAt},
		TestEndAt:   DateTime{testEndAt},
	}
}

func (t *Term) IsInTerm(user *User) bool {
	if user == nil {
		now := flextime.Now()
		return t.StartAt.Before(now) && t.EndAt.After(now)
	}
	now := flextime.Now().Add(user.TimeDifference)
	if user.IsSuperUser() {
		return t.TestStartAt.Before(now) && t.TestEndAt.After(now)
	}
	return t.StartAt.Before(now) && t.EndAt.After(now)
}
