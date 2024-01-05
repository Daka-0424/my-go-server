package model

import "github.com/Daka-0424/my-go-server/pkg/domain/entity"

type VcPlatformProduct struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	PaidPoint   uint   `json:"paid_point"`
	FreePoint   uint   `json:"free_point"`
	ProductId   string `json:"product_id"`
}

type VcPlatformProductList struct {
	Products []VcPlatformProduct `json:"products"`
}

func NewVcPlatformProduct(product *entity.VcPlatformProduct) *VcPlatformProduct {
	return &VcPlatformProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		PaidPoint:   product.PaidPoint,
		FreePoint:   product.FreePoint,
		ProductId:   product.ProductId,
	}
}

func NewVcPlatformProductList(products []entity.VcPlatformProduct) *VcPlatformProductList {
	var productList []VcPlatformProduct
	for _, product := range products {
		productList = append(productList, *NewVcPlatformProduct(&product))
	}

	return &VcPlatformProductList{
		Products: productList,
	}
}
