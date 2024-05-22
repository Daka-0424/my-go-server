package appstore

import (
	"context"
	"fmt"

	"github.com/Daka-0424/my-go-server/pkg/domain/appstore"
	"github.com/awa/go-iap/appstore/api"
)

type AppStore struct {
	client        *api.StoreClient
	sandboxClient *api.StoreClient
}

func NewAppStoreAPI(productionClientConfig, sandboxClientConfig *api.StoreConfig) appstore.IAppStore {
	return &AppStore{
		client:        api.NewStoreClient(productionClientConfig),
		sandboxClient: api.NewStoreClient(sandboxClientConfig),
	}
}

func (a AppStore) GetTrancaction(ctx context.Context, transactionID string) (*api.JWSTransaction, error) {
	return a.getTransaction(ctx, transactionID)
}

func (a AppStore) getTransactionInfo(ctx context.Context, transactionId string, useSandbox bool) (*api.JWSTransaction, error) {
	var client *api.StoreClient
	if useSandbox {
		client = a.sandboxClient
	} else {
		client = a.client
	}

	response, err := client.GetTransactionInfo(ctx, transactionId)
	if err != nil {
		return nil, err
	}

	transaction, err := client.ParseSignedTransaction(response.SignedTransactionInfo)
	if err != nil {
		return nil, err
	}

	if transaction.TransactionID != transactionId {
		return nil, fmt.Errorf("transaction id mismatch")
	}

	return transaction, nil
}

// TransactionIDが見つからない場合にsandboxで再度検索する
func (a AppStore) getTransaction(ctx context.Context, transactionID string) (*api.JWSTransaction, error) {
	response, err := a.getTransactionInfo(ctx, transactionID, false)
	if err != nil {
		if apiError, ok := err.(*api.Error); ok {
			switch apiError.ErrorCode() {
			case 4040010: // TransactionIdNotFound Error
				return a.getTransactionInfo(ctx, transactionID, true)
			}
		}
		return nil, err
	}
	return response, nil
}
