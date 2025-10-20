package api

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/request"
	"github.com/gin-gonic/gin"
)

type PlaystoreController struct {
	controllerBase
	playstoreUsecase usecase.IPlaystore
}

func NewPlaystoreController(
	pu usecase.IPlaystore,
	cfg *config.Config,
	lc *language.Localizer,
) *PlaystoreController {
	return &PlaystoreController{
		playstoreUsecase: pu,
		controllerBase: controllerBase{
			cfg:       cfg,
			localizer: lc,
		},
	}
}

func (ctl *PlaystoreController) Billing(ctx *gin.Context) {
	clime, apperr := ctl.getClaims(ctx)
	if apperr != nil {
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	var req request.PlayStoreBilling
	if apperr := formatter.ShouldBind(ctx, &req, ctl.localizer); apperr != nil {
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	id, apperr := clime.GetUserId(ctx, ctl.localizer)
	if apperr != nil {
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	payment, err := ctl.playstoreUsecase.PlaystoreBilling(ctx, id, req)
	if err != nil {
		apperr := ctl.toAppError(ctx, err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, 200, payment)
}
