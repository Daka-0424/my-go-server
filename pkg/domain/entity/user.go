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
	DisplayCode    string              `gorm:"display_code;size:16"` // 表示用のコード
	UUID           string              `gorm:"uuid;index;size:255"`
	Name           string              `gorm:"name;index;size:255"`
	UserKind       uint                `gorm:"user_kind"`
	TimeDifference time.Duration       `gorm:"time_difference"`
	ClientVersion  string              `gorm:"client_version;index;size:255"`
	Device         string              `gorm:"device;index;size:255"`
	PlatformNumber uint                `gorm:"platform_number"`
	Setting        UserSetting         `gorm:"foreignkey:UserID"`
	Vc             UserSummaryRelation `gorm:"foreignkey:UserID"`
	LoginState     UserLoginState      `gorm:"foreignkey:UserID"`

	gormAuxiliary
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

func (u *User) SetDisplayCode(code string) {
	u.DisplayCode = code
	u.updateColumn("display_code")
}

func (u *User) UpdateDevice(ClientVersion, Device string, platformNumber uint) bool {
	update := false
	if u.ClientVersion != ClientVersion {
		u.ClientVersion = ClientVersion
		update = true
		u.updateColumn("client_version")
	}
	if u.Device != Device {
		u.Device = Device
		update = true
		u.updateColumn("device")
	}
	if u.PlatformNumber != platformNumber {
		u.PlatformNumber = platformNumber
		update = true
		u.updateColumn("platform_number")
	}

	return update
}
