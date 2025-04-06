package infra

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Memory interface {
	Redis() RedisClient
}

type RedisClient interface {
	Ping(ctx context.Context) error
	Close() error
}

type MemoryCE struct {
	redisClient RedisClient
}

func NewMemoryCE(redisClient RedisClient) *MemoryCE {
	return &MemoryCE{
		redisClient: redisClient,
	}
}

func (m *MemoryCE) Redis() RedisClient {
	return m.redisClient
}
