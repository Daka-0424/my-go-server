package repository

import "context"

type ITransaction interface {
	DoInTx(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
