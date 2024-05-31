package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IPaymentAppstoreToken interface {
	CreateOrUpdate(ctx context.Context, appstoreToken *entity.PaymentAppstoreToken) error
	ExistsPaymentAppstoreToken(ctx context.Context, transactionID string) (bool, error)
}
