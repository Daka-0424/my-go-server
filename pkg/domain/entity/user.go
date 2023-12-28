package entity

import "gorm.io/gorm"

const (
	NewUser   = iota // 新規ユーザ
	Player           // 通常ユーザ
	SuperUser        // スーパーユーザ
	Banned           // アカウント停止
	Deleted          // アカウント削除
)

type User struct {
	gorm.Model
	Uuid     string `gorm:"index;size:255"`
	Name     string `gorm:"index;size:255"`
	UserKind uint
}

func (u *User) IsEmpty() bool {
	return u.ID == 0
}

func (u *User) IsSuperUser() bool {
	return u.UserKind == SuperUser
}
