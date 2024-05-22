package playstore

import (
	"context"

	"google.golang.org/api/androidpublisher/v3"
)

type IGooglePlayStoreClient interface {
	VerifyProduct(ctx context.Context, packageName, productID, token string) (*androidpublisher.ProductPurchase, error)
	AcknowledgeProduct(ctx context.Context, packageName, productID, token, developerPayload string) error
	//ConsumeProduct(ctx context.Context, packageName, productID, token string) error
}

type IGooglePlayStore interface {
	VerifySignature(ctx context.Context, receipt []byte, signature string) (bool, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
