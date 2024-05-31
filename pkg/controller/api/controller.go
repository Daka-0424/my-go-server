package controller

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/middleware"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	HeaderPlatform   = "APP_PLATFORM"
	HeaderDevice     = "APP_DEVICE"
	HeaderAppVersion = "APP_VERSION"
)

const (
	PlatformUnknown = iota
	PlatformAndroid
	PlatformIOS
	PlatformWebgl
	PlatformWindows
)

type controllerBase struct {
	cfg       *config.Config
	localizer *i18n.Localizer
}

func (ctl *controllerBase) getPlatform(ctx *gin.Context) (string, uint) {
	platform := ctx.GetHeader(HeaderPlatform)
	switch platform {
	case "Android":
		return platform, PlatformAndroid
	case "iOS":
		return platform, PlatformIOS
	case "WebGL":
		return platform, PlatformWebgl
	case "Windows":
		return platform, PlatformWindows
	default:
		return platform, PlatformUnknown
	}
}

func (ctl *controllerBase) getAppDevice(ctx *gin.Context) string {
	return ctx.GetHeader(HeaderDevice)
}

func (ctl *controllerBase) getAppVersion(ctx *gin.Context) string {
	return ctx.GetHeader(HeaderAppVersion)
}

func (ctl *controllerBase) toAppError(err error) *model.AppError {
	switch apperr := err.(type) {
	case *model.AppError:
		return apperr
	default:
		cf := &i18n.LocalizeConfig{MessageID: model.E9999}
		return model.NewErrInternalServerError(model.E9999, ctl.localizer.MustLocalize(cf))
	}
}

func (ctl *controllerBase) getClaims(ctx *gin.Context) (*middleware.Claims, *model.AppError) {
	claims, ok := ctx.Get("claims")
	if !ok {
		cf := &i18n.LocalizeConfig{MessageID: model.E0101}
		return nil, model.NewErrInternalServerError(model.E0101, ctl.localizer.MustLocalize(cf))
	}
	return claims.(*middleware.Claims), nil
}
