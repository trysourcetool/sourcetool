package infra

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/trysourcetool/sourcetool/backend/config"
)

type RedisClient interface {
	Ping(ctx context.Context) error
	Close() error
}

type RedisClientCE struct {
	client *redis.Client
}

func NewRedisClientCE() *RedisClientCE {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
	})
	
	return &RedisClientCE{
		client: client,
	}
}

func (r *RedisClientCE) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisClientCE) Close() error {
	return r.client.Close()
}
