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

type SessionController struct {
	controllerBase
	sessionUsecase usecase.Session
}

func NewSessionController(
	su usecase.Session,
	cfg *config.Config,
	lc *i18n.Localizer,
) *SessionController {
	return &SessionController{
		controllerBase: controllerBase{cfg: cfg, localizer: lc},
		sessionUsecase: su,
	}
}

func (ctl *SessionController) CreateSession(ctx *gin.Context) {
	var req CreateSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E0101}
		apperr := model.NewErrUnprocessable(model.E0101, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	session, err := ctl.sessionUsecase.CreateSession(ctx, req.UserId, req.Uuid)
	if err != nil {
		apperr := ctl.toAppError(err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, http.StatusOK, session)
}

type CreateSessionRequest struct {
	UserId uint   `json:"user_id"`
	Uuid   string `json:"uuid"`
}
