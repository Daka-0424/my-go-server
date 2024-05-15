package playstore

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/store/playstore"
)

type playstoreFactory struct {
	cfg *config.Config
}

func NewPlaystoreFactory(cfg *config.Config) playstore.PlaystoreFactory {
	return &playstoreFactory{
		cfg: cfg,
	}
}

func (a *playstoreFactory) CreatePlaystoreClient(ctx context.Context) (playstore.GooglePlayStoreClient, error) {
	googleApplicationCredentials := a.cfg.GooglePlay.GoogleApplicationCredentials
	return NewGooglePlayStoreClient(googleApplicationCredentials)
}

func (a *playstoreFactory) CreatePlaystore(ctx context.Context) (playstore.GooglePlayStore, error) {
	base64EncodedPublicKey := a.cfg.GooglePlay.Base64EncodedPublicKey
	return NewGooglePlayStore(base64EncodedPublicKey)
}
