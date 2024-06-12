package controller

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AppstoreController struct {
	controllerBase
	appstoreUsecase usecase.IAppstore
}

func NewAppstoreController(
	au usecase.IAppstore,
	cfg *config.Config,
	lc *i18n.Localizer,
) *AppstoreController {
	return &AppstoreController{
		appstoreUsecase: au,
		controllerBase: controllerBase{
			cfg:       cfg,
			localizer: lc,
		},
	}
}

func (ctl *AppstoreController) Billing(ctx *gin.Context) {
	clime, apperr := ctl.getClaims(ctx)
	if apperr != nil {
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	var req billingAppstoreRequest
	if err := formatter.ShouldBind(ctx, &req); err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9901}
		apperr := model.NewErrUnprocessable(model.E9901, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	id, err := clime.GetUserId()
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E1001}
		apperr = model.NewErrInternalServerError(model.E1001, ctl.localizer.MustLocalize(c))
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	payment, err := ctl.appstoreUsecase.AppstoreBilling(ctx, id, req.TransactionID, req.PurchaseItemID)
	if err != nil {
		apperr := ctl.toAppError(err)
		formatter.Respond(ctx, apperr.StatusCode, gin.H{"error": apperr})
		return
	}

	formatter.Respond(ctx, 200, payment)
}

type billingAppstoreRequest struct {
	TransactionID  string `json:"transaction_id"`
	PurchaseItemID uint   `json:"purchase_item_id"`
}
