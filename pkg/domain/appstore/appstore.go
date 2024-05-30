package appstore

import (
	"context"

	"github.com/awa/go-iap/appstore/api"
)

type IAppStore interface {
	GetTransaction(ctx context.Context, transactionID string) (*api.JWSTransaction, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
