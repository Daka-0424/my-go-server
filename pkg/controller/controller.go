package controller

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/middleware"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type controllerBase struct {
	cfg       *config.Config
	localizer *i18n.Localizer
}

func (c *controllerBase) toAppError(err error) *model.AppError {
	switch apperr := err.(type) {
	case *model.AppError:
		return apperr
	default:
		cf := &i18n.LocalizeConfig{MessageID: model.E9999}
		return model.NewErrInternalServerError(model.E9999, c.localizer.MustLocalize(cf))
	}
}

func (c *controllerBase) getClaims(ctx *gin.Context) (*middleware.Claims, *model.AppError) {
	claims, ok := ctx.Get("claims")
	if !ok {
		cf := &i18n.LocalizeConfig{MessageID: model.E0101}
		return nil, model.NewErrInternalServerError(model.E0101, c.localizer.MustLocalize(cf))
	}
	return claims.(*middleware.Claims), nil
}
