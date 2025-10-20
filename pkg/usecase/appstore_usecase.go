package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/domain/appstore"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/request"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/awa/go-iap/appstore/api"
)

type IAppstore interface {
	AppstoreBilling(ctx context.Context, userID uint, req request.AppStoreBilling) (*response.ReceiptResult, error)
}

type appstoreUsecase struct {
	appstoreFactory               appstore.IAppstoreFactory
	appstoreRepository            repository.IPaymentAppstoreToken
	userSummaryRelationRepository repository.IUserSummaryRelation
	userPointSummaryRepository    repository.IUserPointSummary
	platformProductRepository     repository.ISeed[entity.PlatformProduct]
	transaction                   repository.ITransaction
	earnedPointService            service.IEarnedPoint
	localizer                     *language.Localizer
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
	localizer *language.Localizer,
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

func (usecase *appstoreUsecase) AppstoreBilling(ctx context.Context, userID uint, req request.AppStoreBilling) (*response.ReceiptResult, error) {
	kpiLogger, err := usecase.kpiLoggerFactory.Create(ctx)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, "", nil))
	}

	vc, err := usecase.userSummaryRelationRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E0002, usecase.localizer.MustLocalize(model.E0002, "", nil))
	}

	// 購入情報のProductIDからPlatformProductを取得する
	platformProduct, err := usecase.platformProductRepository.GetByID(ctx, req.PurchaseItemID)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E3001, usecase.localizer.MustLocalize(model.E3001, "", nil))
	}

	appstore, err := usecase.appstoreFactory.Create(ctx)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, "", nil))
	}

	tx, err := appstore.GetTransaction(ctx, req.TransactionID)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E9002, usecase.localizer.MustLocalize(model.E9002, "", nil))
	}

	if tx.TransactionID != req.TransactionID {
		return nil, response.NewErrBadRequest(model.E9005, usecase.localizer.MustLocalize(model.E9005, "", nil))
	}

	if tx.Type != api.Consumable {
		return nil, response.NewErrBadRequest(model.E9006, usecase.localizer.MustLocalize(model.E9006, "", nil))
	}

	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {

		// レシートが存在するかチェックする
		existsAppstoreToken, err := usecase.appstoreRepository.ExistsPaymentAppstoreToken(ctx, tx.TransactionID)
		if err != nil {
			return nil, response.NewErrBadRequest(model.E3001, usecase.localizer.MustLocalize(model.E3001, "", nil))
		}

		if existsAppstoreToken {
			return nil, response.NewErrBadRequest(model.E9007, usecase.localizer.MustLocalize(model.E9007, "", nil))
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
				return nil, response.NewErrUnprocessable(model.E0001, usecase.localizer.MustLocalize(model.E0001, "", nil))
			}
			return nil, response.NewErrUnprocessable(model.E9004, usecase.localizer.MustLocalize(model.E9004, "", nil))
		}

		if platformProduct.PaidPoint > 0 {
			earnedPaidPoint, err := usecase.earnedPointService.Payout(ctx, vc, platformProduct.PaidPoint, entity.GemKindPaid, platformProduct, "by-receipt", appToken.CreatedAt)
			if err != nil {
				return nil, response.NewErrUnprocessable(model.E9102, usecase.localizer.MustLocalize(model.E9102, "", nil))
			}

			appToken.EarnedPointID = earnedPaidPoint.ID
		}

		if platformProduct.FreePoint > 0 {
			_, err := usecase.earnedPointService.Payout(ctx, vc, platformProduct.FreePoint, entity.GemKindFree, platformProduct, "by-receipt", appToken.CreatedAt)
			if err != nil {
				return nil, response.NewErrUnprocessable(model.E9102, usecase.localizer.MustLocalize(model.E9102, "", nil))
			}
		}

		if err := usecase.userPointSummaryRepository.BulkUpdate(ctx, vc.PointSummaries()); err != nil {
			return nil, response.NewErrUnprocessable(model.E9103, usecase.localizer.MustLocalize(model.E9103, "", nil))
		}

		if err := usecase.appstoreRepository.CreateOrUpdate(ctx, appToken); err != nil {
			return nil, response.NewErrUnprocessable(model.E0001, usecase.localizer.MustLocalize(model.E0001, "", nil))
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

		return response.NewReceiptResult(vc, platformProduct), nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*response.ReceiptResult), nil
}
