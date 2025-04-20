package port

import "context"

// Message represents a message received from a subscription.
type Message struct {
	ID      string
	Payload []byte
}

// PubSub defines the interface for publish/subscribe operations.
type PubSub interface {
	// Publish sends a message to the specified channel.
	Publish(ctx context.Context, channel, id string, payload []byte) error
	// Subscribe listens for messages on the specified channel.
	// It returns a channel that receives messages. The channel will be closed
	// when the context is canceled or an error occurs during subscription.
	Subscribe(ctx context.Context, channel string) (<-chan *Message, error)
	// Close terminates the connection.
	Close() error
}
