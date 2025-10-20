package api

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/request"
	"github.com/gin-gonic/gin"
)

type SessionController struct {
	controllerBase
	sessionUsecase usecase.ISession
}

func NewSessionController(
	su usecase.ISession,
	cfg *config.Config,
	lc *language.Localizer,
) *SessionController {
	return &SessionController{
		controllerBase: controllerBase{cfg: cfg, localizer: lc},
		sessionUsecase: su,
	}
}

func (ctl *SessionController) CreateSession(ctx *gin.Context) {
	var req request.Session
	if apperr := formatter.ShouldBind(ctx, &req, ctl.localizer); apperr != nil {
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	device := ctl.getAppDevice(ctx)
	appVersion := ctl.getAppVersion(ctx)
	_, platformNumber := ctl.getPlatform(ctx)

	session, err := ctl.sessionUsecase.CreateSession(ctx, device, appVersion, platformNumber, req)
	if err != nil {
		apperr := ctl.toAppError(ctx, err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, http.StatusOK, session)
}
