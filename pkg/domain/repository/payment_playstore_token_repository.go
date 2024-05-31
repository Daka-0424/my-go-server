package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IPaymentPlaystoreToken interface {
	CreateOrUpdate(ctx context.Context, playstoreToken *entity.PaymentPlaystoreToken) error
	ExistsPaymentPlaystoreToken(ctx context.Context, orderID string) (bool, error)
}
