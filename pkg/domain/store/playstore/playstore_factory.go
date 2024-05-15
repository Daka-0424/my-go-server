package playstore

import "context"

type PlaystoreFactory interface {
	CreatePlaystore(ctx context.Context) (GooglePlayStore, error)
	CreatePlaystoreClient(ctx context.Context) (GooglePlayStoreClient, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
