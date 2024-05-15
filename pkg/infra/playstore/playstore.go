package playstore

import (
	"context"

	"github.com/awa/go-iap/playstore"
	"google.golang.org/api/androidpublisher/v3"
)

type GooglePlayStore struct {
	Base64EncodedPublicKey string
}

type GooglePlayStoreClient struct {
	client *playstore.Client
}

func NewGooglePlayStore(base64EncodedPublicKey string) (*GooglePlayStore, error) {
	return &GooglePlayStore{
		Base64EncodedPublicKey: base64EncodedPublicKey,
	}, nil
}

func NewGooglePlayStoreClient(googleApplicationCredentials string) (*GooglePlayStoreClient, error) {
	client, err := playstore.New([]byte(googleApplicationCredentials))
	if err != nil {
		return &GooglePlayStoreClient{}, err
	}
	return &GooglePlayStoreClient{
		client: client,
	}, nil
}

func (g *GooglePlayStore) VerifySignature(ctx context.Context, receipt []byte, signature string) (bool, error) {
	return playstore.VerifySignature(g.Base64EncodedPublicKey, receipt, signature)
}

func (g *GooglePlayStoreClient) VerifyProduct(ctx context.Context, packageName, productID, token string) (*androidpublisher.ProductPurchase, error) {
	ret, err := g.client.VerifyProduct(ctx, packageName, productID, token)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (g *GooglePlayStoreClient) AcknowledgeProduct(ctx context.Context, packageName, productID, token, developerPayload string) error {
	return g.client.AcknowledgeProduct(ctx, packageName, productID, token, developerPayload)
}
