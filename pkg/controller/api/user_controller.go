package api

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	controllerBase
	registrationUsecase usecase.IUser
}

func NewUserController(
	ru usecase.IUser,
	cfg *config.Config,
	lc *language.Localizer) *UserController {
	return &UserController{
		controllerBase:      controllerBase{cfg: cfg, localizer: lc},
		registrationUsecase: ru,
	}
}

// Registration
//
// @Summary      User Registration
// @Schemes
// @Description  Register a new user
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.User
// @Failure      400  {object}  response.AppError
// @Router /api/registration [post]
func (ctl *UserController) Registration(ctx *gin.Context) {
	device := ctl.getAppDevice(ctx)
	appVersion := ctl.getAppVersion(ctx)
	_, platformNumber := ctl.getPlatform(ctx)

	registration, err := ctl.registrationUsecase.Registration(ctx, device, appVersion, platformNumber)
	if err != nil {
		apperr := ctl.toAppError(ctx, err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, http.StatusOK, registration)
}
