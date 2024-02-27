package usecase

import (
	"context"
	"errors"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type User interface {
	Registration(ctx context.Context, uuid, device, appVersion string, platform uint) (*model.User, error)
}

type userUsercase struct {
	cfg            *config.Config
	localizer      *i18n.Localizer
	transaction    repository.Transaction
	userRepository repository.User
	userService    service.User
}

func NewUserUsecase(
	cfg *config.Config,
	lc *i18n.Localizer,
	transaction repository.Transaction,
	userRepository repository.User,
	userService service.User) User {
	return &userUsercase{
		cfg:            cfg,
		localizer:      lc,
		transaction:    transaction,
		userRepository: userRepository,
		userService:    userService,
	}
}

func (u *userUsercase) Registration(ctx context.Context, uuid, device, appVersion string, platform uint) (*model.User, error) {
	value, err := u.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		user, err := u.userRepository.FindByUuid(ctx, uuid)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			c := &i18n.LocalizeConfig{MessageID: model.E0101}
			return nil, model.NewErrUnprocessable(model.E0101, u.localizer.MustLocalize(c))
		}

		if user != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E0102}
			return nil, model.NewErrUnprocessable(model.E0102, u.localizer.MustLocalize(c))
		}

		user, err = u.userService.CreateUser(ctx, uuid, device, appVersion, platform)
		if err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E0103}
			return nil, model.NewErrUnprocessable(model.E0103, u.localizer.MustLocalize(c))
		}

		return model.NewUser(user), nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*model.User), nil
}
