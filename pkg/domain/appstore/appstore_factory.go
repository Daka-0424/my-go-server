package appstore

import "context"

type IAppstoreFactory interface {
	Create(ctx context.Context) (IAppStore, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
