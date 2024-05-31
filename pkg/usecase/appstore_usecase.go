package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/appstore"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/awa/go-iap/appstore/api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type IAppstore interface {
	AppstoreBilling(ctx context.Context, userID uint, transactionID string, purchaseItemID uint) (*model.ReceiptResult, error)
}

type appstoreUsecase struct {
	appstoreFactory               appstore.IAppstoreFactory
	appstoreRepository            repository.IPaymentAppstoreToken
	userSummaryRelationRepository repository.IUserSummaryRelation
	userPointSummaryRepository    repository.IUserPointSummary
	platformProductRepository     repository.ISeed[entity.PlatformProduct]
	transaction                   repository.ITransaction
	earnedPointService            service.IEarnedPoint
	localizer                     *i18n.Localizer
	kpiLoggerFactory              logger.IKpiLoggerFactory
}

func NewAppstoreUsecase(
	appstoreFactory appstore.IAppstoreFactory,
	par repository.IPaymentAppstoreToken,
	usrr repository.IUserSummaryRelation,
	upsr repository.IUserPointSummary,
	pp repository.ISeed[entity.PlatformProduct],
	tx repository.ITransaction,
	eps service.IEarnedPoint,
	localizer *i18n.Localizer,
	kpiLoggerFactory logger.IKpiLoggerFactory,
) IAppstore {
	return &appstoreUsecase{
		appstoreFactory:               appstoreFactory,
		appstoreRepository:            par,
		userSummaryRelationRepository: usrr,
		userPointSummaryRepository:    upsr,
		platformProductRepository:     pp,
		transaction:                   tx,
		earnedPointService:            eps,
		localizer:                     localizer,
		kpiLoggerFactory:              kpiLoggerFactory,
	}
}

func (usecase *appstoreUsecase) AppstoreBilling(ctx context.Context, userID uint, transactionID string, purchaseItemID uint) (*model.ReceiptResult, error) {
	kpiLogger, err := usecase.kpiLoggerFactory.Create(ctx)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9999}
		return nil, model.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(c))
	}

	vc, err := usecase.userSummaryRelationRepository.FindByUserID(ctx, userID)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E0002}
		return nil, model.NewErrBadRequest(model.E0002, usecase.localizer.MustLocalize(c))
	}

	// 購入情報のProductIDからPlatformProductを取得する
	platformProduct, err := usecase.platformProductRepository.GetByID(ctx, purchaseItemID)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E3001}
		return nil, model.NewErrBadRequest(model.E3001, usecase.localizer.MustLocalize(c))
	}

	appstore, err := usecase.appstoreFactory.Create(ctx)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9999}
		return nil, model.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(c))
	}

	tx, err := appstore.GetTransaction(ctx, transactionID)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9002}
		return nil, model.NewErrBadRequest(model.E9002, usecase.localizer.MustLocalize(c))
	}

	if tx.TransactionID != transactionID {
		c := &i18n.LocalizeConfig{MessageID: model.E9005}
		return nil, model.NewErrBadRequest(model.E9005, usecase.localizer.MustLocalize(c))
	}

	if tx.Type != api.Consumable {
		c := &i18n.LocalizeConfig{MessageID: model.E9006}
		return nil, model.NewErrBadRequest(model.E9006, usecase.localizer.MustLocalize(c))
	}

	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {

		// レシートが存在するかチェックする
		existsAppstoreToken, err := usecase.appstoreRepository.ExistsPaymentAppstoreToken(ctx, tx.TransactionID)
		if err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E3001}
			return nil, model.NewErrBadRequest(model.E3001, usecase.localizer.MustLocalize(c))
		}

		if existsAppstoreToken {
			c := &i18n.LocalizeConfig{MessageID: model.E9007}
			return nil, model.NewErrBadRequest(model.E9007, usecase.localizer.MustLocalize(c))
		}

		appToken := entity.NewPaymentAppstoreToken(
			tx.TransactionID,
			tx.AppAccountToken,
			tx.BundleID,
			tx.Currency,
			string(tx.Environment),
			tx.ProductID,
			uint(tx.Price),
			uint(tx.PurchaseDate),
			uint(tx.Quantity),
			uint(tx.RevocationDate),
			vc.UserID,
			platformProduct,
		)

		if tx.RevocationDate > 0 {
			if err := usecase.appstoreRepository.CreateOrUpdate(ctx, appToken); err != nil {
				c := &i18n.LocalizeConfig{MessageID: model.E0001}
				return nil, model.NewErrUnprocessable(model.E0001, usecase.localizer.MustLocalize(c))
			}
			c := &i18n.LocalizeConfig{MessageID: model.E9004}
			return nil, model.NewErrUnprocessable(model.E9004, usecase.localizer.MustLocalize(c))
		}

		if platformProduct.PaidPoint > 0 {
			earnedPaidPoint, err := usecase.earnedPointService.Payout(ctx, vc, platformProduct.PaidPoint, entity.GemKindPaid, platformProduct, "by-receipt", appToken.CreatedAt)
			if err != nil {
				c := &i18n.LocalizeConfig{MessageID: model.E9102}
				return nil, model.NewErrUnprocessable(model.E9102, usecase.localizer.MustLocalize(c))
			}

			appToken.EarnedPointID = earnedPaidPoint.ID
		}

		if platformProduct.FreePoint > 0 {
			_, err := usecase.earnedPointService.Payout(ctx, vc, platformProduct.FreePoint, entity.GemKindFree, platformProduct, "by-receipt", appToken.CreatedAt)
			if err != nil {
				c := &i18n.LocalizeConfig{MessageID: model.E9102}
				return nil, model.NewErrUnprocessable(model.E9102, usecase.localizer.MustLocalize(c))
			}
		}

		if err := usecase.userPointSummaryRepository.BulkUpdate(ctx, vc.PointSummaries()); err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E9103}
			return nil, model.NewErrUnprocessable(model.E9103, usecase.localizer.MustLocalize(c))
		}

		if err := usecase.appstoreRepository.CreateOrUpdate(ctx, appToken); err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E0001}
			return nil, model.NewErrUnprocessable(model.E0001, usecase.localizer.MustLocalize(c))
		}

		kpiDate := map[string]interface{}{
			"user_id":             vc.UserID,
			"platform_product_id": platformProduct.ID,
			"product_name":        platformProduct.Name,
			"point":               platformProduct.FreePoint,
			"point_rate":          platformProduct.UnitCost(),
			"price":               platformProduct.Price,
			"currency_key":        "JPY",
		}
		kpiLogger.LogEvent(logger.KpiLogPayment, kpiDate)
		kpiLogger.Flush()

		return model.NewReceiptResult(vc, platformProduct), nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*model.ReceiptResult), nil
}
