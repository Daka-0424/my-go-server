package controller

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type PlaystoreController struct {
	controllerBase
	playstoreUsecase usecase.IPlaystore
}

func NewPlaystoreController(
	pu usecase.IPlaystore,
	cfg *config.Config,
	lc *i18n.Localizer,
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

	var req billingPlaystoreRequest
	if err := formatter.ShouldBind(ctx, &req); err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9901}
		apperr := model.NewErrUnprocessable(model.E9901, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	id, err := clime.GetUserId()
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E0101}
		apperr = model.NewErrInternalServerError(model.E0101, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	payment, err := ctl.playstoreUsecase.PlaystoreBilling(ctx, id, req.PurchaseItemID, req.Receipt, req.Signature)
	if err != nil {
		apperr := ctl.toAppError(err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, 200, payment)
}

type billingPlaystoreRequest struct {
	PurchaseItemID uint   `json:"purchaseItemId"`
	Receipt        string `json:"receipt"`
	Signature      string `json:"signature"`
}
