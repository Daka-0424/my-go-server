package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	RedirectEndPoint = "/oauth/google/users/%s"
)

type IGoogle interface {
}

type googleUsecase struct {
	localizer          *i18n.Localizer
	oauthGoogleService service.IOauthGoogle
	transaction        repository.ITransaction
	userRepository     repository.IUser
	kpiLoggerFactory   logger.IKpiLoggerFactory
}

func NewGoogleUsecase(
	lc *i18n.Localizer,
	ogs service.IOauthGoogle,
	rt repository.ITransaction,
	ur repository.IUser,
	klfl logger.IKpiLoggerFactory,
) IGoogle {
	return &googleUsecase{
		localizer:          lc,
		oauthGoogleService: ogs,
		transaction:        rt,
		userRepository:     ur,
		kpiLoggerFactory:   klfl,
	}
}

func (usecase *googleUsecase) OauthGoogleURL(ctx context.Context, userID uint) model.GoogleLoginURL {
	url := usecase.oauthGoogleService.OauthGoogleURL(ctx, oauthRedirectURL(userID))
	return model.NewGoogleLoginURL(url)
}

func (usecase *googleUsecase) OauthGoogle(ctx context.Context, code string, userID uint) error {
	_, err := usecase.oauthGoogleService.OauthGoogle(ctx, code, oauthRedirectURL(userID))
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E1101}
		return model.NewErrBadRequest(model.E1101, usecase.localizer.MustLocalize(c))
	}
	return nil
}

func oauthRedirectURL(userID uint) string {
	return fmt.Sprintf(RedirectEndPoint, strconv.Itoa(int(userID)))
}
