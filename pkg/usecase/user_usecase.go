package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type IUser interface {
	Registration(ctx context.Context, uuid, device, clientVersion string, platformNumber uint) (*model.User, error)
}

type userUsercase struct {
	cfg                  *config.Config
	localizer            *i18n.Localizer
	transaction          repository.ITransaction
	userRepository       repository.IUser
	loginStateRepository repository.IUserLoginState
	userService          service.IUser
	vcService            service.IVc
	kpiLoggerFactory     logger.IKpiLoggerFactory
}

func NewUserUsecase(
	cfg *config.Config,
	lc *i18n.Localizer,
	transaction repository.ITransaction,
	userRepository repository.IUser,
	loginStateRepository repository.IUserLoginState,
	userService service.IUser,
	vcService service.IVc,
	kpiLoggerFactory logger.IKpiLoggerFactory,
) IUser {
	return &userUsercase{
		cfg:                  cfg,
		localizer:            lc,
		transaction:          transaction,
		userRepository:       userRepository,
		loginStateRepository: loginStateRepository,
		userService:          userService,
		vcService:            vcService,
		kpiLoggerFactory:     kpiLoggerFactory,
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

		fns := []func(ctx context.Context, user *entity.User) error{
			usecase.createUserLoginState,
		}

		for _, fn := range fns {
			if err := fn(ctx, user); err != nil {
				return nil, err
			}
		}

		kpiLogger, err := usecase.kpiLoggerFactory.Create(ctx)
		if err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E9999}
			return nil, model.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(c))
		}

		loginDate := map[string]interface{}{
			"user_id": user.ID,
		}
		kpiLogger.LogEvent(logger.KpiLogInstall, loginDate)
		kpiLogger.Flush()

		return model.NewUser(user), nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*model.User), nil
}

func (usecase *userUsercase) createUserLoginState(ctx context.Context, user *entity.User) error {
	loginState := entity.NewUserLoginState(user.ID)

	if err := usecase.loginStateRepository.CreateOrUpdate(ctx, loginState); err != nil {
		return err
	}

	return nil
}
