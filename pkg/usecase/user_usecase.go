package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type IUser interface {
	Registration(ctx context.Context, uuid, device, clientVersion string, platformNumber uint) (*model.User, error)
}

type userUsercase struct {
	cfg            *config.Config
	localizer      *i18n.Localizer
	transaction    repository.ITransaction
	userRepository repository.IUser
	userService    service.IUser
	vcService      service.IVc
}

func NewUserUsecase(
	cfg *config.Config,
	lc *i18n.Localizer,
	transaction repository.ITransaction,
	userRepository repository.IUser,
	userService service.IUser,
	vcService service.IVc,
) IUser {
	return &userUsercase{
		cfg:            cfg,
		localizer:      lc,
		transaction:    transaction,
		userRepository: userRepository,
		userService:    userService,
		vcService:      vcService,
	}
}

func (usecase *userUsercase) Registration(ctx context.Context, uuid, device, clientVersion string, platformNumber uint) (*model.User, error) {
	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		if uuid == "" {
			c := &i18n.LocalizeConfig{MessageID: model.E9901}
			return nil, model.NewErrUnprocessable(model.E9901, usecase.localizer.MustLocalize(c))
		}

		// もし、uuidで検索して、userが存在していたら、エラーを返す
		exists, err := usecase.userRepository.ExistsUser(ctx, uuid)
		if err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E0002}
			return nil, model.NewErrUnprocessable(model.E0002, usecase.localizer.MustLocalize(c))
		}
		if exists {
			c := &i18n.LocalizeConfig{MessageID: model.E0106}
			return nil, model.NewErrUnprocessable(model.E0106, usecase.localizer.MustLocalize(c))
		}

		// なかったら、新規登録する
		user, err := usecase.userService.Register(ctx, uuid, device, clientVersion, platformNumber)
		if err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E0103}
			return nil, model.NewErrUnprocessable(model.E0103, usecase.localizer.MustLocalize(c))
		}
		// VCのセットアップ
		if err := usecase.vcService.SetupVc(ctx, user); err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E0103}
			return nil, model.NewErrUnprocessable(model.E0103, usecase.localizer.MustLocalize(c))
		}

		return model.NewUser(user), nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*model.User), nil
}
