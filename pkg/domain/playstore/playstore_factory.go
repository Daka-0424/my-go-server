package playstore

import "context"

type IPlaystoreFactory interface {
	CreatePlaystore(ctx context.Context) (IGooglePlayStore, error)
	CreatePlaystoreClient(ctx context.Context) (IGooglePlayStoreClient, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
