package infra

import (
	"context"
)

type RedisClient interface {
	Ping(ctx context.Context) error
	Close() error
}
