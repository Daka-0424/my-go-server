package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	NewUser   = iota // 新規ユーザ
	Player           // 通常ユーザ
	SuperUser        // スーパーユーザ
	Banned           // アカウント停止
	Deleted          // アカウント削除
)

type User struct {
	gorm.Model
	Uuid           string `gorm:"index;size:255"`
	Name           string `gorm:"index;size:255"`
	UserKind       uint
	TimeDifference time.Duration
	AppVersion     string `gorm:"index;size:255"`
	Device         string `gorm:"index;size:255"`
	PlatformNumber uint
}

func (u *User) IsEmpty() bool {
	return u.ID == 0
}

func (u *User) IsSuperUser() bool {
	return u.UserKind == SuperUser
}

func (u *User) DisplayCode() string {
	first := rune((u.CreatedAt.Year() - 2000 + 45) % 256)
	second := rune((u.CreatedAt.Month() + 67) % 256)
	return string(first) + u.Uuid + string(second)
}
