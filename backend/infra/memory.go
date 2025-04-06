package infra

import (
	"context"
)

type Memory interface {
	Redis() RedisClient
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
