package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/playstore"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/request"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
)

type IPlaystore interface {
	PlaystoreBilling(ctx context.Context, userID uint, req request.PlayStoreBilling) (*response.ReceiptResult, error)
}

type playstoreUsecase struct {
	playstoreFactory              playstore.IPlaystoreFactory
	playstoreRepository           repository.IPaymentPlaystoreToken
	userSummaryRelationRepository repository.IUserSummaryRelation
	userPointSummaryRepository    repository.IUserPointSummary
	platformProductRepository     repository.ISeed[entity.PlatformProduct]
	transaction                   repository.ITransaction
	earnedPointService            service.IEarnedPoint
	localizer                     *language.Localizer
	kpiLoggerFactory              logger.IKpiLoggerFactory
}

type ReceiptJson struct {
	OrderID            string `json:"orderId"`
	PackageName        string `json:"packageName"`
	ProductID          string `json:"productId"`
	PurchaseTimeMillis string `json:"purchaseTimeMillis"`
	PurchaseState      int    `json:"purchaseState"`
	PurchaseToken      string `json:"purchaseToken"`
	AutoRenewing       bool   `json:"autoRenewing"`
	Acknowlodged       bool   `json:"acknowledged"`
	DeveloperPayload   string `json:"developerPayload"`
}

func NewPlaystoreUsecase(
	pfp playstore.IPlaystoreFactory,
	pptr repository.IPaymentPlaystoreToken,
	usrr repository.IUserSummaryRelation,
	upsr repository.IUserPointSummary,
	ppr repository.ISeed[entity.PlatformProduct],
	tr repository.ITransaction,
	eps service.IEarnedPoint,
	lc *language.Localizer,
	klfl logger.IKpiLoggerFactory,
) IPlaystore {
	return &playstoreUsecase{
		playstoreFactory:              pfp,
		playstoreRepository:           pptr,
		userSummaryRelationRepository: usrr,
		userPointSummaryRepository:    upsr,
		platformProductRepository:     ppr,
		transaction:                   tr,
		earnedPointService:            eps,
		localizer:                     lc,
		kpiLoggerFactory:              klfl,
	}
}

func (usecase *playstoreUsecase) PlaystoreBilling(ctx context.Context, userID uint, req request.PlayStoreBilling) (*response.ReceiptResult, error) {
	kpiLogger, err := usecase.kpiLoggerFactory.Create(ctx)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	vc, err := usecase.userSummaryRelationRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E0002, usecase.localizer.MustLocalize(model.E0002, language.LanguageJapanese, nil))
	}

	platformProduct, err := usecase.platformProductRepository.GetByID(ctx, req.PurchaseItemID)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E3001, usecase.localizer.MustLocalize(model.E3001, language.LanguageJapanese, nil))
	}

	playstore, err := usecase.playstoreFactory.CreatePlaystore(ctx)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	// receiptのデコード
	decoded, err := base64.StdEncoding.DecodeString(req.Receipt)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E3001, usecase.localizer.MustLocalize(model.E3001, language.LanguageJapanese, nil))
	}

	// receiptの検証
	verifySignature, err := playstore.VerifySignature(ctx, decoded, req.Signature)
	if err != nil {
		return nil, response.NewErrUnprocessable(model.E9002, usecase.localizer.MustLocalize(model.E9002, language.LanguageJapanese, nil))
	}

	if !verifySignature {
		return nil, response.NewErrUnprocessable(model.E9001, usecase.localizer.MustLocalize(model.E9001, language.LanguageJapanese, nil))
	}

	// receiptのJSON変換
	receiptJSON, err := parseGoogleReceiptJSON(decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse receipt json: %w", err)
	}

	playstoreClient, err := usecase.playstoreFactory.CreatePlaystoreClient(ctx)
	if err != nil {
		return nil, response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	// receiptから購入情報を所得
	verifyProduct, err := playstoreClient.VerifyProduct(ctx, receiptJSON.PackageName, receiptJSON.ProductID, receiptJSON.PurchaseToken)
	if err != nil {
		return nil, response.NewErrUnprocessable(model.E9006, usecase.localizer.MustLocalize(model.E9006, language.LanguageJapanese, nil))
	}

	if verifyProduct.PurchaseState != 0 {
		return nil, response.NewErrUnprocessable(model.E9004, usecase.localizer.MustLocalize(model.E9004, language.LanguageJapanese, nil))
	}

	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {

		existsPlaystoreToken, err := usecase.playstoreRepository.ExistsPaymentPlaystoreToken(ctx, receiptJSON.OrderID)
		if err != nil {
			return nil, response.NewErrBadRequest(model.E9002, usecase.localizer.MustLocalize(model.E9002, language.LanguageJapanese, nil))
		}

		if existsPlaystoreToken {
			return nil, response.NewErrBadRequest(model.E9007, usecase.localizer.MustLocalize(model.E9007, language.LanguageJapanese, nil))
		}

		playstoreToken := entity.NewPaymentPlaystoreToken(
			verifyProduct.OrderId,
			receiptJSON.PackageName,
			verifyProduct.ProductId,
			verifyProduct.PurchaseState,
			verifyProduct.PurchaseToken,
			verifyProduct.PurchaseTimeMillis,
			verifyProduct.Quantity,
			verifyProduct.RegionCode,
			verifyProduct.ConsumptionState,
			verifyProduct.Kind,
			verifyProduct.DeveloperPayload,
			verifyProduct.AcknowledgementState,
			verifyProduct.ObfuscatedExternalAccountId,
			verifyProduct.ObfuscatedExternalProfileId,
			platformProduct.ID,
			req.Signature,
			entity.GetPurchaseType(verifyProduct.PurchaseType),
			vc.UserID,
		)

		if platformProduct.PaidPoint > 0 {
			earnedPaidPoint, err := usecase.earnedPointService.Payout(ctx, vc, platformProduct.PaidPoint, entity.GemKindPaid, platformProduct, "by-receipt", playstoreToken.CreatedAt)
			if err != nil {
				return nil, err
			}

			playstoreToken.EarnedPointID = earnedPaidPoint.ID
		}

		if platformProduct.FreePoint > 0 {
			_, err := usecase.earnedPointService.Payout(ctx, vc, platformProduct.FreePoint, entity.GemKindFree, platformProduct, "by-receipt", playstoreToken.CreatedAt)
			if err != nil {
				return nil, err
			}
		}

		if err := usecase.userPointSummaryRepository.BulkUpdate(ctx, vc.PointSummaries()); err != nil {
			return nil, response.NewErrUnprocessable(model.E9103, usecase.localizer.MustLocalize(model.E9103, language.LanguageJapanese, nil))
		}

		if err := usecase.playstoreRepository.CreateOrUpdate(ctx, playstoreToken); err != nil {
			return nil, response.NewErrUnprocessable(model.E0001, usecase.localizer.MustLocalize(model.E0001, language.LanguageJapanese, nil))
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

func parseGoogleReceiptJSON(decodedReceip []byte) (*ReceiptJson, error) {
	var receiptJSON ReceiptJson
	if err := json.Unmarshal(decodedReceip, &receiptJSON); err != nil {
		return nil, err
	}

	return &receiptJSON, nil
}
