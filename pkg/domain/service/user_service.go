package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

const (
	USER_DEFAULT_NAME = "NewUser"
	CodeKey           = 81
)

type UserData struct {
	UUID           string
	Device         string
	ClientVersion  string
	PlatformNumber uint
	LanguageCode   string
}

type IUser interface {
	Register(ctx context.Context, userData *UserData) (*entity.User, error)
}

type userService struct {
	userRepository           repository.IUser
	userLoginStateRepository repository.IUserLoginState
	userSettingRepository    repository.IUserSetting
}

func NewUserService(
	userRepository repository.IUser,
	userLoginStateRepository repository.IUserLoginState,
	userSettingRepository repository.IUserSetting,
) IUser {
	return &userService{
		userRepository:           userRepository,
		userLoginStateRepository: userLoginStateRepository,
		userSettingRepository:    userSettingRepository,
	}
}

func (service *userService) Register(ctx context.Context, userData *UserData) (*entity.User, error) {
	user, err := service.userRepository.CreateUser(ctx, userData.UUID, USER_DEFAULT_NAME, userData.ClientVersion, userData.Device, userData.PlatformNumber)
	if err != nil {
		return nil, err
	}

	// DisplayCodeを作成
	user.SetDisplayCode(service.createDisplayCode(user))
	if err := service.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	if err := service.createUserSetting(ctx, user, userData.LanguageCode); err != nil {
		return nil, err
	}

	fns := []func(ctx context.Context, user *entity.User) error{
		service.createUserLoginState,
	}

	for _, fn := range fns {
		if err := fn(ctx, user); err != nil {
			return nil, err
		}
	}

	// その後、userを返す
	return user, nil
}

func (service *userService) createUserSetting(ctx context.Context, user *entity.User, languageCode string) error {
	userSetting := entity.NewUserSetting(user.ID, languageCode)
	if err := service.userSettingRepository.CreateOrUpdate(ctx, userSetting); err != nil {
		return err
	}

	user.Setting = *userSetting
	return nil
}

func (service *userService) createUserLoginState(ctx context.Context, user *entity.User) error {
	loginState := entity.NewUserLoginState(user.ID)

	if err := service.userLoginStateRepository.CreateOrUpdate(ctx, loginState); err != nil {
		return err
	}

	user.LoginState = *loginState
	return nil
}

func (service *userService) createDisplayCode(user *entity.User) string {
	first := rune((user.CreatedAt.Year() - 2020 + 45) % 256)
	second := rune((int(user.CreatedAt.Month()) + 67) % 256)
	code := hash(user.ID)

	return fmt.Sprintf("%c%c%02d", first, second, code)
}

func hash(userID uint) int {
	multiplier := rand.Intn(CodeKey)
	h := int(userID)
	h = h*multiplier + CodeKey
	return h
}
