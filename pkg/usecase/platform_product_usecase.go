package usecase

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type IPlatformProduct interface {
	ListPlatformProducts(ctx context.Context) (*model.PlatformProductList, error)
	FindPlatformNumber(ctx context.Context, platformNumber uint) (*model.PlatformProductList, error)
}

type platformProductUsecase struct {
	cfg                           *config.Config
	localizer                     *i18n.Localizer
	seedPlatformProductRepository repository.ISeed[entity.PlatformProduct]
}

func NewPlatformProductUsecase(
	cfg *config.Config,
	lc *i18n.Localizer,
	seedPlatformProduct repository.ISeed[entity.PlatformProduct]) IPlatformProduct {
	return &platformProductUsecase{
		cfg:                           cfg,
		localizer:                     lc,
		seedPlatformProductRepository: seedPlatformProduct,
	}
}

func (usecase *platformProductUsecase) ListPlatformProducts(ctx context.Context) (*model.PlatformProductList, error) {
	products, err := usecase.seedPlatformProductRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	list := model.NewPlatformProductList(products)

	return list, nil
}

func (usecase *platformProductUsecase) FindPlatformNumber(ctx context.Context, platformNumber uint) (*model.PlatformProductList, error) {
	products, err := usecase.seedPlatformProductRepository.Where(ctx, entity.PlatformProduct{PlatformNumber: platformNumber})
	if err != nil {
		return nil, err
	}

	list := model.NewPlatformProductList(products)

	return list, nil
}
