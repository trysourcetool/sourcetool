package ws

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/config"
	redisv1 "github.com/trysourcetool/sourcetool/proto/go/redis/v1"
)

type redisClient struct {
	client *redis.Client
}

func newRedisClient() (*redisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisClient{
		client: client,
	}, nil
}

func (r *redisClient) Publish(ctx context.Context, channel string, id string, payload []byte) error {
	message := &redisv1.RedisMessage{
		Id:      id,
		Payload: payload,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal redis message: %w", err)
	}

	return r.client.Publish(ctx, channel, data).Err()
}

func (r *redisClient) Subscribe(ctx context.Context, channel string) (<-chan *redis.Message, error) {
	pubsub := r.client.Subscribe(ctx, channel)
	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel: %w", err)
	}

	return pubsub.Channel(), nil
}

func (r *redisClient) Close() error {
	return r.client.Close()
}
