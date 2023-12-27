package util

import (
	"context"
	"time"

	"github.com/Daka-0424/my-go-server/pkg/domain/util"
	"github.com/go-redsync/redsync/v4"
)

type mutexFactoryUtil struct {
	redsync *redsync.Redsync
}

func NewLockFactoryUtil(redsync *redsync.Redsync) util.MutexFactory {
	return &mutexFactoryUtil{redsync: redsync}
}

func (s *mutexFactoryUtil) Create(key string, ttl time.Duration) util.Mutex {
	mutex := s.redsync.NewMutex(key, redsync.WithExpiry(ttl), redsync.WithTries(100), redsync.WithRetryDelay(10*time.Millisecond))
	return NewMutexUtil(mutex)
}

type mutexUtil struct {
	mutex *redsync.Mutex
}

func NewMutexUtil(mutex *redsync.Mutex) util.Mutex {
	return &mutexUtil{mutex: mutex}
}

func (s *mutexUtil) Lock(ctx context.Context) error {
	return s.mutex.LockContext(ctx)
}

func (s *mutexUtil) Unlock(ctx context.Context) (bool, error) {
	return s.mutex.UnlockContext(ctx)
}
