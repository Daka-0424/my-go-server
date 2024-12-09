package controller

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type PlatformProductController struct {
	controllerBase
	platformProductUsecase usecase.IPlatformProduct
}

func NewPlatformProductController(pu usecase.IPlatformProduct, cfg *config.Config, lc *i18n.Localizer) *PlatformProductController {
	return &PlatformProductController{
		controllerBase:         controllerBase{cfg: cfg, localizer: lc},
		platformProductUsecase: pu,
	}
}

func (ctl *PlatformProductController) ListPlatformProduct(ctx *gin.Context) {
	_, apperr := ctl.getClaims(ctx)
	if apperr != nil {
		formatter.Respond(ctx, ctl.cfg, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	_, platformNumber := ctl.getPlatform(ctx)
	products, err := ctl.platformProductUsecase.FindPlatformNumber(ctx, platformNumber)
	if err != nil {
		apperr := ctl.toAppError(err)
		formatter.Respond(ctx, ctl.cfg, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, ctl.cfg, http.StatusOK, products)
}
