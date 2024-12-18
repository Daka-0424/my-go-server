package entity

import (
	"gorm.io/gorm"
)

type AdminRoleType string

// TODO:一旦3Typeに分けたが、Masterの登録をどうするかは別途考える
const (
	AdminRoleTypeMaster AdminRoleType = "master" // 管理画面の全機能を利用可能
	AdminRoleTypeBasic  AdminRoleType = "basic"  // 管理者の登録/削除/編集以外の機能を利用可能
	AdminRoleTypeGuest  AdminRoleType = "guest"  // 管理画面の一部機能を利用可能
)

type Admin struct {
	gorm.Model
	Email    string        `gorm:"email;index;size:255;not null"`
	Password string        `gorm:"password;not null"`
	RoleType AdminRoleType `gorm:"role_type; not null"`
}

func NewAdmin(email, password string, roleType AdminRoleType) *Admin {
	return &Admin{
		Email:    email,
		Password: password,
		RoleType: roleType,
	}
}
