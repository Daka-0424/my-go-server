package api

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type UserController struct {
	controllerBase
	registrationUsecase usecase.IUser
}

func NewUserController(
	ru usecase.IUser,
	cfg *config.Config,
	lc *i18n.Localizer) *UserController {
	return &UserController{
		controllerBase:      controllerBase{cfg: cfg, localizer: lc},
		registrationUsecase: ru,
	}
}

func (ctl *UserController) Registration(ctx *gin.Context) {
	var rew CreateRegistrationRequest
	if err := formatter.ShouldBind(ctx, ctl.cfg, &rew); err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E0101}
		apperr := model.NewErrUnprocessable(model.E0101, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, ctl.cfg, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	device := ctl.getAppDevice(ctx)
	appVersion := ctl.getAppVersion(ctx)
	_, platformNumber := ctl.getPlatform(ctx)

	registration, err := ctl.registrationUsecase.Registration(ctx, rew.Uuid, device, appVersion, platformNumber)
	if err != nil {
		apperr := ctl.toAppError(err)
		formatter.Respond(ctx, ctl.cfg, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, ctl.cfg, http.StatusOK, registration)
}

type CreateRegistrationRequest struct {
	Uuid string `json:"uuid"`
}
