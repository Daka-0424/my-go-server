package model

import "github.com/Daka-0424/my-go-server/pkg/domain/entity"

type Registration struct {
	User User `json:"user"`
}

func NewRegistration(u *entity.User) *Registration {
	return &Registration{
		User: *NewUser(u),
	}
}
