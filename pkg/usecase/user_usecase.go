package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/google/uuid"
)

type IUser interface {
	Registration(ctx context.Context, device, clientVersion string, platformNumber uint) (*response.User, error)
}

type userUsercase struct {
	cfg                  *config.Config
	localizer            *language.Localizer
	transaction          repository.ITransaction
	userRepository       repository.IUser
	loginStateRepository repository.IUserLoginState
	userService          service.IUser
	vcService            service.IVc
	kpiLoggerFactory     logger.IKpiLoggerFactory
}

func NewUserUsecase(
	cfg *config.Config,
	lc *language.Localizer,
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

func (usecase *userUsercase) Registration(ctx context.Context, device, clientVersion string, platformNumber uint) (*response.User, error) {
	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		uuid := uuid.NewString()

		// なかったら、新規登録する
		user, err := usecase.userService.Register(ctx, uuid, device, clientVersion, platformNumber)
		if err != nil {
			return nil, response.NewErrUnprocessable(model.E0103, usecase.localizer.MustLocalize(model.E0103, language.LanguageJapanese, nil))
		}
		// VCのセットアップ
		if err := usecase.vcService.SetupVc(ctx, user); err != nil {
			return nil, response.NewErrUnprocessable(model.E0103, usecase.localizer.MustLocalize(model.E0103, language.LanguageJapanese, nil))
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
			return nil, response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
		}

		kpiDate := map[string]interface{}{
			"user_id": user.ID,
		}
		kpiLogger.LogEvent(logger.KpiLogInstall, kpiDate)
		kpiLogger.Flush()

		return response.NewUser(user), nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*response.User), nil
}

func (usecase *userUsercase) createUserLoginState(ctx context.Context, user *entity.User) error {
	loginState := entity.NewUserLoginState(user.ID)

	if err := usecase.loginStateRepository.CreateOrUpdate(ctx, loginState); err != nil {
		return err
	}

	return nil
}
