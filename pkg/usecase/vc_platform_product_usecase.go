package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type VcPlatformProduct interface {
	ListVcPlatformProducts(ctx context.Context) (*model.VcPlatformProductList, error)
	FindPlatformNumber(ctx context.Context, platformNumber uint) (*model.VcPlatformProductList, error)
}

type vcPlatformProductUsecase struct {
	cfg                             *config.Config
	localizer                       *i18n.Localizer
	seedVcPlatformProductRepository repository.Seed[entity.VcPlatformProduct]
}

func NewVcPlatformProductUsecase(
	cfg *config.Config,
	lc *i18n.Localizer,
	seedVcPlatformProduct repository.Seed[entity.VcPlatformProduct]) VcPlatformProduct {
	return &vcPlatformProductUsecase{
		cfg:                             cfg,
		localizer:                       lc,
		seedVcPlatformProductRepository: seedVcPlatformProduct,
	}
}

func (u *vcPlatformProductUsecase) ListVcPlatformProducts(ctx context.Context) (*model.VcPlatformProductList, error) {
	products, err := u.seedVcPlatformProductRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	list := model.NewVcPlatformProductList(products)

	return list, nil
}

func (u *vcPlatformProductUsecase) FindPlatformNumber(ctx context.Context, platformNumber uint) (*model.VcPlatformProductList, error) {
	products, err := u.seedVcPlatformProductRepository.Where(ctx, entity.VcPlatformProduct{PlatformNumber: platformNumber})
	if err != nil {
		return nil, err
	}

	list := model.NewVcPlatformProductList(products)

	return list, nil
}
