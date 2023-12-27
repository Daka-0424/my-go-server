package util

import (
	"context"
	"time"
)

type MutexFactory interface {
	Create(key string, ttl time.Duration) Mutex
}

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) (bool, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
