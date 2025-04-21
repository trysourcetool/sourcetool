package internal

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

func OpenRedis() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return redisClient, nil
}
