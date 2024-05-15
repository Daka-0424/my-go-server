package appstore

import "context"

type AppstoreFactory interface {
	Create(ctx context.Context) (AppStore, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
