package mock

import (
	"sync"

	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

type client struct {
	handler    websocket.MessageHandlerFunc
	messages   []*websocketv1.Message
	messagesMu sync.RWMutex
	done       chan error
}

func NewClient() *client {
	return &client{
		messages: make([]*websocketv1.Message, 0),
		done:     make(chan error, 1),
	}
}

func (c *client) RegisterHandler(handler websocket.MessageHandlerFunc) {
	c.handler = handler
}

func (c *client) Enqueue(id string, payload proto.Message) {
	msg, err := websocket.NewMessage(id, payload)
	if err != nil {
		return
	}

	c.messagesMu.Lock()
	c.messages = append(c.messages, msg)
	c.messagesMu.Unlock()

	if c.handler != nil {
		c.handler(msg)
	}
}

func (c *client) EnqueueWithResponse(id string, payload proto.Message) (*websocketv1.Message, error) {
	msg, err := websocket.NewMessage(id, payload)
	if err != nil {
		return nil, err
	}

	c.messagesMu.Lock()
	c.messages = append(c.messages, msg)
	c.messagesMu.Unlock()

	return msg, nil
}

func (c *client) Messages() []*websocketv1.Message {
	c.messagesMu.RLock()
	defer c.messagesMu.RUnlock()
	return c.messages
}

func (c *client) Close() error {
	c.done <- nil
	return nil
}

func (c *client) Wait() error {
	return <-c.done
}
