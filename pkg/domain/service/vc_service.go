package service

import (
	"context"
	"fmt"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type IVc interface {
	SetupVc(ctx context.Context, user *entity.User) error
}

type vcService struct {
	userRepository                repository.IUser
	userSummaryRelationRepository repository.IUserSummaryRelation
	localizer                     *i18n.Localizer
}

func NewVcService(
	ur repository.IUser,
	usrr repository.IUserSummaryRelation,
	localizer *i18n.Localizer,
) IVc {
	return &vcService{
		userRepository:                ur,
		userSummaryRelationRepository: usrr,
		localizer:                     localizer,
	}
}

func (service *vcService) SetupVc(ctx context.Context, user *entity.User) error {
	user, vc, _ := service.userRepository.FindUserWithVc(ctx, user.ID)
	// ユーザーのVCが有効でない場合のみセットアップを行う

	if vc.ID == 0 {
		// 有効でない場合はVCを作成
		user.Vc = *entity.NewUserSummaryRelation(user.ID, user.PlatformNumber)
	}

	userVcFreePointSummary, _ := service.userRepository.FindUserPointSummary(ctx, user.ID, user.PlatformNumber, 0)
	userVcPaidPointSummary, _ := service.userRepository.FindUserPointSummary(ctx, user.ID, user.PlatformNumber, 1)

	// 有償間はPF間で別管理なので別PFのVCがない場合のみ作成
	if userVcPaidPointSummary == nil {
		otherPlatformVc, _ := service.userRepository.FindOtherPlatformVc(ctx, user.ID, user.PlatformNumber)
		fmt.Println(otherPlatformVc)

		if otherPlatformVc != nil {
			user.Vc.PaidPointSummary = otherPlatformVc.PaidPointSummary
		} else {
			user.Vc.PaidPointSummary = *entity.NewUserPointSummary(user.ID, 1)
		}
	}

	// 無償間はPF間共通なのでFreePointSummary既存の場合はそれを使う
	if userVcFreePointSummary == nil {

		// 既存のFreePointSummaryがある場合はそれを使う
		fmt.Println(user.ID)
		freePointSummary, err := service.userRepository.FirstOrCreateFreePointSummary(ctx, user.ID, 0)
		if err != nil {
			return err
		}
		user.Vc.FreePointSummary = *freePointSummary
	}

	// VCを保存
	if err := service.userRepository.CreateUserSummaryRelation(ctx, &user.Vc); err != nil {
		return err
	}

	return nil
}

func (service *vcService) IsUserVcValid(userSummary *entity.UserSummaryRelation) bool {
	// ユーザーのVCが有効であるかのロジックを実装
	//user.Vcがnilの場合は無効
	//user.Vc.FreePointSummaryがnilの場合は無効
	//user.Vc.PaidPointSummaryがnilの場合は無効

	return userSummary != nil || userSummary.FreePointSummaryID != 0 || userSummary.PaidPointSummaryID != 0
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
