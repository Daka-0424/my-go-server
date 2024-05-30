package service

import (
	"context"
	"time"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type IEarnedPoint interface {
	Payout(ctx context.Context, vc *entity.UserSummaryRelation, amount uint, earnedType uint, vcPlatformProduct *entity.PlatformProduct, earnSource string, earnedAt time.Time) (*entity.EarnedPoint, error)
	PayoutImitation(ctx context.Context, vc *entity.UserSummaryRelation, withFreePoint uint, platformProduct *entity.PlatformProduct, imitationPoint *entity.ImitationPoint, earnedAt time.Time) (*entity.EarnedPoint, *entity.EarnedPoint, error)
}

type earnedPointService struct {
	earnedPointRepository      repository.IEarnedPoint
	userPointSummaryRepository repository.IUserPointSummary
	appleReceiptRepository     repository.IPaymentAppstoreToken
}

func NewEarnedPointService(
	epr repository.IEarnedPoint,
	upsr repository.IUserPointSummary,
	patr repository.IPaymentAppstoreToken,
) IEarnedPoint {
	return &earnedPointService{
		earnedPointRepository:      epr,
		userPointSummaryRepository: upsr,
		appleReceiptRepository:     patr,
	}
}

func (service *earnedPointService) Payout(ctx context.Context, vc *entity.UserSummaryRelation, amount uint, earnedType uint, vcPlatformProduct *entity.PlatformProduct, earnSource string, earnedAt time.Time) (*entity.EarnedPoint, error) {
	if amount <= 0 {
		return nil, nil
	}
	switch earnedType {
	case entity.GemKindFree:
		earnedPoint := entity.NewCreateFreeEarnedPoint(vc, amount, vcPlatformProduct, earnSource, nil)
		if err := service.earnedPointRepository.CreateOrUpdate(ctx, earnedPoint); err != nil {
			return nil, err
		}

		vc.FreePointSummary.EarnPoint += amount
		vc.FreePointSummary.UpdateBalancePoint()
		if err := service.userPointSummaryRepository.Update(ctx, &vc.FreePointSummary); err != nil {
			return nil, err
		}
		return earnedPoint, nil
	case entity.GemKindPaid:
		earnedPoint := entity.NewCreatePaidEarnedPoint(vc, amount, vcPlatformProduct, earnSource, nil, time.Now())
		if err := service.earnedPointRepository.CreateOrUpdate(ctx, earnedPoint); err != nil {
			return nil, err
		}

		vc.PaidPointSummary.EarnPoint += amount
		vc.PaidPointSummary.UpdateBalancePoint()
		if err := service.userPointSummaryRepository.Update(ctx, &vc.PaidPointSummary); err != nil {
			return nil, err
		}
		return earnedPoint, nil
	}

	return nil, nil
}

func (service *earnedPointService) PayoutImitation(ctx context.Context, vc *entity.UserSummaryRelation, withFreePoint uint, platformProduct *entity.PlatformProduct, imitationPoint *entity.ImitationPoint, earnedAt time.Time) (*entity.EarnedPoint, *entity.EarnedPoint, error) {
	var earnedPoints []entity.EarnedPoint
	var paidEarnedPoint, freeEarnedPoint *entity.EarnedPoint
	if platformProduct.PaidPoint > 0 {
		paidEarnedPoint = entity.NewCreatePaidEarnedPoint(vc, platformProduct.PaidPoint, platformProduct, "by-imitation", imitationPoint, earnedAt)
		earnedPoints = append(earnedPoints, *paidEarnedPoint)
	}
	if withFreePoint > 0 || platformProduct.FreePoint > 0 {
		freeEarnedPoint = entity.NewCreateFreeEarnedPoint(vc, platformProduct.FreePoint, platformProduct, "by-imitation", imitationPoint)
		earnedPoints = append(earnedPoints, *freeEarnedPoint)
	}

	if err := service.earnedPointRepository.BulkCreate(ctx, earnedPoints); err != nil {
		return nil, nil, err
	}

	return paidEarnedPoint, freeEarnedPoint, nil
}
