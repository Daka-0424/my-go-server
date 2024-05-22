package util

import (
	"context"
	"time"
)

type IMutexFactory interface {
	Create(key string, ttl time.Duration) IMutex
}

type IMutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) (bool, error)
}

//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../../mock/$GOPACKAGE/$GOFILE
