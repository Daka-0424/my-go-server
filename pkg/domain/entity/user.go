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
	DisplayCode    string              `gorm:"display_code", "size:16"` // 表示用のコード
	UUID           string              `gorm:"uuid", "index;size:255"`
	Name           string              `gorm:"name", "index;size:255"`
	UserKind       uint                `gorm:"user_kind"`
	TimeDifference time.Duration       `gorm:"time_difference"`
	ClientVersion  string              `gorm:"client_version", "index;size:255"`
	Device         string              `gorm:"device", "index;size:255"`
	PlatformNumber uint                `gorm:"platform_number"`
	Vc             UserSummaryRelation `gorm:"foreignkey:UserID"`
	LoginState     UserLoginState      `gorm:"foreignkey:UserID"`
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

func (u *User) UpdateUserKind(kind uint) {
	u.UserKind = kind
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
