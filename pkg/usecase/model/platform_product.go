package model

import "github.com/Daka-0424/my-go-server/pkg/domain/entity"

type PlatformProduct struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	PaidPoint   uint   `json:"paid_point"`
	FreePoint   uint   `json:"free_point"`
	ProductId   string `json:"product_id"`
}

type PlatformProductList struct {
	Products []PlatformProduct `json:"products"`
}

func NewPlatformProduct(product *entity.PlatformProduct) *PlatformProduct {
	return &PlatformProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		PaidPoint:   product.PaidPoint,
		FreePoint:   product.FreePoint,
		ProductId:   product.ProductId,
	}
}

func NewPlatformProductList(products []entity.PlatformProduct) *PlatformProductList {
	var productList []PlatformProduct
	for _, product := range products {
		productList = append(productList, *NewPlatformProduct(&product))
	}

	return &PlatformProductList{
		Products: productList,
	}
}
