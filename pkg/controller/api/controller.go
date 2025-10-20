package api

import (
	"strconv"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/middleware"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/gin-gonic/gin"
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
	localizer *language.Localizer
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

func (c *controllerBase) toAppError(ctx *gin.Context, err error) *response.AppError {
	switch apperr := err.(type) {
	case *response.AppError:
		return apperr
	default:
		lang := ctx.Request.Header.Get("Accept-Language")
		return response.NewErrInternalServerError(model.E9999, c.localizer.MustLocalize(model.E9999, lang, nil))
	}
}

func (c *controllerBase) requestError(ctx *gin.Context) *response.AppError {
	lang := ctx.Request.Header.Get("Accept-Language")
	return response.NewErrBadRequest(model.E9901, c.localizer.MustLocalize(model.E9901, lang, nil))
}

func (c *controllerBase) getClaims(ctx *gin.Context) (*middleware.Claims, *response.AppError) {
	climes, ok := ctx.Get("claims")
	if !ok {
		lang := ctx.Request.Header.Get("Accept-Language")
		err := response.NewErrInternalServerError(model.E0101, c.localizer.MustLocalize(model.E0101, lang, nil))
		return nil, err
	}

	return climes.(*middleware.Claims), nil
}

func (c *controllerBase) parseUint(ctx *gin.Context, param string) (uint, *response.AppError) {
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		lang := ctx.Request.Header.Get("Accept-Language")
		return 0, response.NewErrBadRequest(model.E9901, c.localizer.MustLocalize(model.E9901, lang, nil))
	}
	return uint(id), nil
}

func (c *controllerBase) parseBool(ctx *gin.Context, param string) (bool, *response.AppError) {
	b, err := strconv.ParseBool(param)
	if err != nil {
		lang := ctx.Request.Header.Get("Accept-Language")
		return false, response.NewErrBadRequest(model.E9901, c.localizer.MustLocalize(model.E9901, lang, nil))
	}
	return b, nil
}
