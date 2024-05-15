package appstore

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/appstore"
	"github.com/awa/go-iap/appstore/api"
)

type appstoreFactory struct {
	cfg *config.Config
}

func NewAppStoreFactory(cfg *config.Config) appstore.AppstoreFactory {
	return &appstoreFactory{
		cfg: cfg,
	}
}

func (a *appstoreFactory) Create(ctx context.Context) (appstore.AppStore, error) {
	prdConfig, err := newPrdAppStoreConfig(a.cfg, false)
	if err != nil {
		return nil, err
	}

	sandboxConfig, err := newSandboxAppStoreConfig(a.cfg, true)
	if err != nil {
		return nil, err
	}

	return NewAppStoreAPI(prdConfig, sandboxConfig), nil
}

// 本番用のAppStoreConfigを生成する
func newPrdAppStoreConfig(cfg *config.Config, isSandbox bool) (*api.StoreConfig, error) {
	// TODO: 環境変数から取得するようにする
	return &api.StoreConfig{
		KeyContent: []byte(cfg.KeyContent), // Loads a .p8 certificate
		KeyID:      cfg.KeyID,              // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
		BundleID:   cfg.BundleID,           // Your app’s bundle ID
		Issuer:     cfg.IssuerID,           // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
		Sandbox:    isSandbox,              // default is Production
	}, nil
}

// サンドボックス用のAppStoreConfigを生成する
func newSandboxAppStoreConfig(cfg *config.Config, isSandbox bool) (*api.StoreConfig, error) {
	return &api.StoreConfig{
		KeyContent: []byte(cfg.SandboxKeyContent), // Loads a .p8 certificate
		KeyID:      cfg.SandboxKeyID,              // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
		BundleID:   cfg.SandboxBundleID,           // Your app’s bundle ID
		Issuer:     cfg.SandboxIssuerID,           // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
		Sandbox:    isSandbox,                     // default is Production
	}, nil
}
