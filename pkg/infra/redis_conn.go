package infra

import (
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisConnector struct {
	Client  *redis.Client
	RedSync *redsync.Redsync
}

func NewRedisConnector(cfg *config.Config) *RedisConnector {
	timeout := 3 * time.Second
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.CONN,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		MaxRetries:   3,
	})

	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &RedisConnector{Client: client, RedSync: rs}
}
