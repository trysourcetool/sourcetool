package pubsub

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	redisv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/redis/v1"
)

type Message struct {
	ID      string
	Payload []byte
}

type PubSub struct {
	redisClient *redis.Client
}

func New(client *redis.Client) *PubSub {
	return &PubSub{client}
}

func (c *PubSub) Publish(ctx context.Context, channel, id string, payload []byte) error {
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

func (c *PubSub) Subscribe(ctx context.Context, channel string) (<-chan *Message, error) {
	redisPubSub := c.redisClient.Subscribe(ctx, channel)
	// Wait for confirmation that subscription is created before receiving anything.
	if _, err := redisPubSub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel %q: %w", channel, err)
	}

	// Create a Go channel to forward messages.
	msgChan := make(chan *Message)

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
				msgChan <- &Message{
					ID:      message.Id,
					Payload: message.Payload,
				}
			}
		}
	}()

	return msgChan, nil
}

func (c *PubSub) Close() error {
	return c.redisClient.Close()
}
