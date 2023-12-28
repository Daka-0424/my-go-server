package model

import "github.com/Daka-0424/my-go-server/pkg/domain/entity"

type Session struct {
	User  User   `json:"user"`
	Token string `json:"token"`
	Key   string `json:"key"`
	IV    string `json:"iv"`
}

func NewSession(u *entity.User, token, key, iv string) *Session {
	return &Session{
		User:  *NewUser(u),
		Token: token,
		Key:   key,
		IV:    iv,
	}
}
