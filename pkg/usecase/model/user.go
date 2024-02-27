package model

import "github.com/Daka-0424/my-go-server/pkg/domain/entity"

type User struct {
	ID       uint   `json:"id"`
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	UserKind uint   `json:"user_kind"`
}

func NewUser(u *entity.User) *User {
	return &User{
		ID:       u.ID,
		Uuid:     u.UUID,
		Name:     u.Name,
		UserKind: u.UserKind,
	}
}
