package repository

import (
	"context"
	"time"

	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) repository.ICache {
	return &redisCache{
		client: client,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	bytes, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return bytes, true, nil
}

func (c *redisCache) Set(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	err := c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
