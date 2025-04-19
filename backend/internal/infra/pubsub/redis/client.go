package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/pubsub"
	redisv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/redis/v1"
)

type clientCE struct {
	redisClient *redis.Client
}

// Compile-time check to ensure client implements the PubSub interface.
var _ pubsub.PubSub = (*clientCE)(nil)

// NewClient creates a new Redis client implementing the PubSub interface.
func NewClientCE() (pubsub.PubSub, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &clientCE{
		redisClient: redisClient,
	}, nil
}

func (c *clientCE) Publish(ctx context.Context, channel, id string, payload []byte) error {
	message := &redisv1.RedisMessage{
		Id:      id,
		Payload: payload,
	}
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal redis message: %w", err)
	}

	return c.redisClient.Publish(ctx, channel, data).Err()
}

func (c *clientCE) Subscribe(ctx context.Context, channel string) (<-chan *pubsub.Message, error) {
	redisPubSub := c.redisClient.Subscribe(ctx, channel)
	// Wait for confirmation that subscription is created before receiving anything.
	if _, err := redisPubSub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel %q: %w", channel, err)
	}

	// Create a Go channel to forward messages.
	msgChan := make(chan *pubsub.Message)

	go func() {
		defer close(msgChan)
		defer redisPubSub.Close() // Ensure Redis PubSub is closed when the goroutine exits.

		redisChan := redisPubSub.Channel()

		for {
			select {
			case <-ctx.Done():
				// Context canceled, stop listening.
				return
			case redisMsg, ok := <-redisChan:
				if !ok {
					// Redis channel closed.
					return
				}

				// Unmarshal the received message.
				var message redisv1.RedisMessage
				err := proto.Unmarshal([]byte(redisMsg.Payload), &message)
				if err != nil {
					// Log error or handle it appropriately.
					// For now, we'll just skip this message.
					// Consider adding logging here.
					fmt.Printf("Error unmarshalling message: %v\n", err) // TODO: Replace with proper logging
					continue
				}

				// Forward the message in the application's format.
				msgChan <- &pubsub.Message{
					ID:      message.Id,
					Payload: message.Payload,
				}
			}
		}
	}()

	return msgChan, nil
}

func (c *clientCE) Close() error {
	return c.redisClient.Close()
}
