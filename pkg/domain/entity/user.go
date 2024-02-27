package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	DefaultUserName = "初心者さん"
)

const (
	Tutee     = iota // 新規ユーザ
	Player           // 通常ユーザ
	SuperUser        // スーパーユーザ
	Banned           // アカウント停止
)

type User struct {
	gorm.Model
	DisplayCode    string `json:"display_code" gorm:"size:16"` // 表示用のコード
	UUID           string `gorm:"index;size:255"`
	Name           string `gorm:"index;size:255"`
	UserKind       uint
	TimeDifference time.Duration
	ClientVersion  string `gorm:"index;size:255"`
	Device         string `gorm:"index;size:255"`
	PlatformNumber uint
	LoginState     UserLoginState `json:"user_login_state" gorm:"foreignkey:UserID"`
}

func NewUser(uuid string, name string, clientVersion string, device string, platformNumber uint) *User {
	return &User{
		UUID:           uuid,
		Name:           name,
		ClientVersion:  clientVersion,
		Device:         device,
		PlatformNumber: platformNumber,
	}
}

func (u *User) IsEmpty() bool {
	return u.ID == 0
}

func (u *User) IsSuperUser() bool {
	return u.UserKind == SuperUser
}

func (u *User) UpdateDevice(ClientVersion, Device string, platformNumber uint) bool {
	if u.ClientVersion == ClientVersion && u.Device == Device && u.PlatformNumber == platformNumber {
		return false
	}

	u.ClientVersion = ClientVersion
	u.Device = Device
	u.PlatformNumber = platformNumber
	return true
}
