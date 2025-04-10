package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool-go/internal/logger"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
)

const (
	// Time constraints
	minPingInterval   = 100 * time.Millisecond
	maxPingInterval   = 30 * time.Second
	minReconnectDelay = 100 * time.Millisecond

	// Queue constraints
	minQueueSize = 50
	maxQueueSize = 1000

	// Default values
	defaultPingInterval   = time.Second
	defaultReconnectDelay = time.Second
	defaultQueueSize      = 250
)

func validateConfig(config Config) error {
	// Minimum values
	if config.PingInterval < minPingInterval {
		return fmt.Errorf("ping interval must be at least %v", minPingInterval)
	}
	if config.ReconnectDelay < minReconnectDelay {
		return fmt.Errorf("reconnect delay must be at least %v", minReconnectDelay)
	}
	if config.QueueSize < minQueueSize {
		return fmt.Errorf("queue size must be at least %d", minQueueSize)
	}

	// Maximum values
	if config.PingInterval > maxPingInterval {
		return fmt.Errorf("ping interval must not exceed %v", maxPingInterval)
	}
	if config.QueueSize > maxQueueSize {
		return fmt.Errorf("queue size must not exceed %d", maxQueueSize)
	}

	return nil
}

func setConfigDefaults(config *Config) {
	if config.PingInterval == 0 {
		config.PingInterval = defaultPingInterval
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = defaultReconnectDelay
	}
	if config.QueueSize == 0 {
		config.QueueSize = defaultQueueSize
	}
}

type Client interface {
	RegisterHandler(MessageHandlerFunc)
	Enqueue(string, proto.Message)
	EnqueueWithResponse(string, proto.Message) (*websocketv1.Message, error)
	Close() error
	Wait() error
}

type Config struct {
	URL            string
	APIKey         string
	InstanceID     uuid.UUID
	PingInterval   time.Duration
	ReconnectDelay time.Duration
	QueueSize      int
	OnReconnecting func()
	OnReconnected  func()
}

type client struct {
	config Config
	dialer *websocket.Dialer

	conn   *websocket.Conn
	connMu sync.RWMutex

	messageQueue chan *websocketv1.Message
	done         chan error

	responses map[string]chan *websocketv1.Message
	respMu    sync.RWMutex

	handler   MessageHandlerFunc
	handlerMu sync.RWMutex

	// Goroutine management
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewClient(config Config) (Client, error) {
	// Set defaults for zero values
	setConfigDefaults(&config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	c := &client{
		config:       config,
		messageQueue: make(chan *websocketv1.Message, config.QueueSize),
		done:         make(chan error, 1),
		dialer:       websocket.DefaultDialer,
		responses:    make(map[string]chan *websocketv1.Message),
		stop:         make(chan struct{}),
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	c.wg.Add(3)
	go c.pingPongLoop()
	go c.readMessages()
	go c.sendEnqueuedMessagesLoop()

	return c, nil
}

func (c *client) RegisterHandler(handler MessageHandlerFunc) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	c.handler = handler
}

func (c *client) handleMessage(msg *websocketv1.Message) error {
	// Handle responses
	c.respMu.RLock()
	respChan, exists := c.responses[msg.Id]
	c.respMu.RUnlock()

	if exists {
		respChan <- msg
		return nil
	}

	// Handle calls
	c.handlerMu.RLock()
	handler := c.handler
	c.handlerMu.RUnlock()

	if handler == nil {
		return fmt.Errorf("no message handler registered")
	}

	return handler(msg)
}

func (c *client) connect() error {
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	header.Set("X-Instance-Id", c.config.InstanceID.String())

	conn, _, err := c.dialer.Dial(c.config.URL, header)
	if err != nil {
		return err
	}

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(c.config.PingInterval * 2))
	})

	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

	return nil
}

func (c *client) reconnect() error {
	if c.config.OnReconnecting != nil {
		c.config.OnReconnecting()
	}

	for {
		err := c.connect()
		if err == nil {
			if c.config.OnReconnected != nil {
				c.config.OnReconnected()
			}
			return nil
		}

		logger.Log.Error("reconnection failed, retrying",
			zap.Error(err),
			zap.Duration("delay", c.config.ReconnectDelay))
		time.Sleep(c.config.ReconnectDelay)
	}
}

func (c *client) pingPongLoop() {
	defer c.wg.Done() // Signal completion on exit
	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.connMu.RLock()
			conn := c.conn
			c.connMu.RUnlock()

			if conn == nil {
				return
			}

			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				logger.Log.Error("ping failed", zap.Error(err))
				c.connMu.Lock()
				if c.conn == conn {
					conn.Close()
					c.conn = nil
				}
				c.connMu.Unlock()

				go func() {
					if err := c.reconnect(); err != nil {
						logger.Log.Error("reconnection failed", zap.Error(err))
					}
				}()
				return
			}
		case <-c.stop:
			logger.Log.Debug("pingPongLoop stopping")
			return
		}
	}
}

func (c *client) readMessages() {
	defer c.wg.Done() // Signal completion on exit
	for {
		select {
		case <-c.stop:
			logger.Log.Debug("readMessages stopping")
			return
		default:
			var data []byte
			c.connMu.RLock()
			conn := c.conn
			c.connMu.RUnlock()

			if conn == nil {
				select {
				case <-c.stop:
					logger.Log.Debug("readMessages stopping (conn nil)")
					return
				case <-time.After(100 * time.Millisecond):
					continue
				}
			}

			_, data, err := conn.ReadMessage()
			if err != nil {
				select {
				case <-c.stop:
					logger.Log.Debug("readMessages stopping (read error during shutdown)")
					return
				default:
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						return
					}

					logger.Log.Error("read error", zap.Error(err))
					c.connMu.Lock()
					if c.conn == conn {
						conn.Close()
						c.conn = nil
					}
					c.connMu.Unlock()

					go func() {
						if err := c.reconnect(); err != nil {
							logger.Log.Error("reconnection failed", zap.Error(err))
						}
					}()
					continue
				}
			}

			msg, err := unmarshalMessage(data)
			if err != nil {
				logger.Log.Error("error unmarshaling message", zap.Error(err))
				continue
			}

			if err := c.handleMessage(msg); err != nil {
				logger.Log.Error("error handling message", zap.Error(err))
			}
		}
	}
}

func (c *client) sendEnqueuedMessagesLoop() {
	defer c.wg.Done() // Signal completion on exit

	const batchInterval = 10 * time.Millisecond
	var messageBuffer []*websocketv1.Message

	for {
		select {
		case <-c.stop:
			logger.Log.Debug("sendEnqueuedMessagesLoop stopping")
			c.connMu.RLock()
			conn := c.conn
			c.connMu.RUnlock()
			if conn != nil && len(messageBuffer) > 0 {
				logger.Log.Info("sending remaining messages before shutdown", zap.Int("count", len(messageBuffer)))
				rateLimiter := time.NewTicker(time.Millisecond)
				defer rateLimiter.Stop()
				for _, msg := range messageBuffer {
					_ = c.send(msg)
					<-rateLimiter.C
				}
			}
			return
		case msg, ok := <-c.messageQueue:
			if !ok {
				logger.Log.Debug("sendEnqueuedMessagesLoop stopping (queue closed)")
				c.connMu.RLock()
				conn := c.conn
				c.connMu.RUnlock()
				if conn != nil && len(messageBuffer) > 0 {
					logger.Log.Info("sending remaining messages on queue close", zap.Int("count", len(messageBuffer)))
					rateLimiter := time.NewTicker(time.Millisecond)
					defer rateLimiter.Stop()
					for _, m := range messageBuffer {
						_ = c.send(m)
						<-rateLimiter.C
					}
				}
				return
			}
			messageBuffer = append(messageBuffer, msg)
		default:
			if len(messageBuffer) > 0 {
				c.connMu.RLock()
				conn := c.conn
				c.connMu.RUnlock()

				if conn == nil {
					time.Sleep(time.Second)
					continue
				}

				var remainingMessages []*websocketv1.Message
				for _, msg := range messageBuffer {
					if err := c.send(msg); err != nil {
						remainingMessages = append(remainingMessages, msg)
						logger.Log.Error("error sending message", zap.Error(err))
						break
					}
					time.Sleep(time.Millisecond)
				}

				messageBuffer = remainingMessages

				if len(remainingMessages) == 0 {
					time.Sleep(batchInterval)
				}
			}
		}
	}
}

func (c *client) send(msg *websocketv1.Message) error {
	data, err := marshalMessage(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (c *client) Enqueue(id string, payload proto.Message) {
	msg, err := NewMessage(id, payload)
	if err != nil {
		logger.Log.Error("error creating message", zap.Error(err))
		return
	}
	c.messageQueue <- msg
}

func (c *client) EnqueueWithResponse(id string, payload proto.Message) (*websocketv1.Message, error) {
	respChan := make(chan *websocketv1.Message, 1)
	c.respMu.Lock()
	c.responses[id] = respChan
	c.respMu.Unlock()

	defer func() {
		c.respMu.Lock()
		delete(c.responses, id)
		c.respMu.Unlock()
	}()

	c.Enqueue(id, payload)

	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

func (c *client) Wait() error {
	return <-c.done
}

func (c *client) Close() error {
	c.connMu.Lock()
	if c.conn == nil {
		c.connMu.Unlock()
		return nil // Already closed or never connected
	}

	logger.Log.Debug("closing websocket client")

	// Signal goroutines to stop
	close(c.stop)

	// Close the message queue to unblock sender if waiting
	// Do this *after* signaling stop to allow sender to potentially process buffer on stop signal
	close(c.messageQueue)

	// Send close message (best effort)
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		// Log error only if it's not a standard close error
		logger.Log.Warn("error sending close message", zap.Error(err))
	}

	// Close the underlying connection
	connErr := c.conn.Close()
	if connErr != nil {
		logger.Log.Error("error closing websocket connection", zap.Error(connErr))
	}
	c.conn = nil
	c.connMu.Unlock() // Unlock before waiting

	// Wait for goroutines to finish
	logger.Log.Debug("waiting for goroutines to stop")
	c.wg.Wait()
	logger.Log.Debug("goroutines stopped")

	// Signal that the client is fully closed
	// Use non-blocking send in case Wait() already returned due to previous error
	select {
	case c.done <- nil:
	default:
	}

	return connErr // Return error from closing the connection, if any
}
