package controller

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type RegistrationController struct {
	controllerBase
	registrationUsecase usecase.Registration
}

func NewRegistrationController(
	ru usecase.Registration,
	cfg *config.Config,
	lc *i18n.Localizer) *RegistrationController {
	return &RegistrationController{
		controllerBase:      controllerBase{cfg: cfg, localizer: lc},
		registrationUsecase: ru,
	}
}

func (ctl *RegistrationController) CreateRegistration(ctx *gin.Context) {
	var rew CreateRegistrationRequest
	if err := formatter.ShouldBind(ctx, &rew); err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E0101}
		apperr := model.NewErrUnprocessable(model.E0101, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	registration, err := ctl.registrationUsecase.CreateRegistration(ctx, rew.Uuid, rew.Name)
	if err != nil {
		apperr := ctl.toAppError(err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, http.StatusOK, registration)
}

type CreateRegistrationRequest struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}
