package websocket

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
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
	// Time constraints.
	minPingInterval   = 100 * time.Millisecond
	maxPingInterval   = 30 * time.Second
	minReconnectDelay = 100 * time.Millisecond

	// Queue constraints.
	minQueueSize = 50
	maxQueueSize = 1000

	// Default values.
	defaultPingInterval   = time.Second
	defaultReconnectDelay = time.Second
	defaultQueueSize      = 250

	// Reconnection constants.
	initialReconnectDelay = 500 * time.Millisecond
	maxReconnectDelay     = 30 * time.Second
	// 4s * 2^13 ≈ 1 hour.
	maxReconnectAttempts = 26 // Approximately 1 hour of reconnection attempts

	// Message sending constants.
	maxMessageRetries = 3
	messageRetryDelay = 100 * time.Millisecond
	shutdownTimeout   = 5 * time.Second
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

	// Shutdown state
	shutdownOnce sync.Once
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

	attempt := 0
	lastSuccessTime := time.Now()

	for {
		// Calculate delay with exponential backoff, but cap it
		delay := min(initialReconnectDelay*time.Duration(1<<uint(attempt)), maxReconnectDelay)

		// Add some jitter to prevent thundering herd
		// Use crypto/rand for better randomness
		maxJitter := int64(delay / 4)
		if maxJitter > 0 {
			jitterBig, err := rand.Int(rand.Reader, big.NewInt(maxJitter))
			if err == nil {
				delay += time.Duration(jitterBig.Int64())
			}
		}

		// Log reconnection attempt with more context
		logger.Log.Info("attempting to reconnect",
			zap.Int("attempt", attempt+1),
			zap.Duration("delay", delay),
			zap.String("instance_id", c.config.InstanceID.String()),
			zap.Time("last_success", lastSuccessTime),
			zap.Duration("time_since_last_success", time.Since(lastSuccessTime)))

		// Try to connect
		err := c.connect()
		if err == nil {
			logger.Log.Info("reconnection successful",
				zap.Int("attempts", attempt+1),
				zap.String("instance_id", c.config.InstanceID.String()),
				zap.Duration("total_downtime", time.Since(lastSuccessTime)))

			if c.config.OnReconnected != nil {
				c.config.OnReconnected()
			}
			return nil
		}

		attempt++

		// Check if we should stop trying
		if attempt >= maxReconnectAttempts {
			// If less than an hour has passed since the last successful connection, stop trying
			if time.Since(lastSuccessTime) < time.Hour {
				logger.Log.Error("max reconnection attempts reached within an hour",
					zap.Error(err),
					zap.String("instance_id", c.config.InstanceID.String()),
					zap.Duration("total_downtime", time.Since(lastSuccessTime)))
				return fmt.Errorf("failed to reconnect after %d attempts within an hour: %w", maxReconnectAttempts, err)
			}

			// If more than an hour has passed, continue reconnection attempts
			logger.Log.Warn("continuing reconnection attempts after an hour",
				zap.Error(err),
				zap.String("instance_id", c.config.InstanceID.String()),
				zap.Duration("total_downtime", time.Since(lastSuccessTime)))
			attempt = 0 // Reset counter
			continue
		}

		// Wait before next attempt
		select {
		case <-c.stop:
			logger.Log.Info("reconnection canceled during shutdown",
				zap.String("instance_id", c.config.InstanceID.String()),
				zap.Duration("total_downtime", time.Since(lastSuccessTime)))
			return nil
		case <-time.After(delay):
			continue
		}
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
	defer c.wg.Done()
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
						logger.Log.Info("connection closed normally",
							zap.String("instance_id", c.config.InstanceID.String()))
						return
					}

					logger.Log.Warn("read error, initiating reconnection",
						zap.Error(err),
						zap.String("instance_id", c.config.InstanceID.String()))

					c.connMu.Lock()
					if c.conn == conn {
						conn.Close()
						c.conn = nil
					}
					c.connMu.Unlock()

					// Start reconnection in a separate goroutine
					go func() {
						if err := c.reconnect(); err != nil {
							logger.Log.Error("reconnection failed",
								zap.Error(err),
								zap.String("instance_id", c.config.InstanceID.String()))
						}
					}()
					continue
				}
			}

			msg, err := unmarshalMessage(data)
			if err != nil {
				logger.Log.Error("error unmarshaling message",
					zap.Error(err),
					zap.String("instance_id", c.config.InstanceID.String()))
				continue
			}

			if err := c.handleMessage(msg); err != nil {
				logger.Log.Error("error handling message",
					zap.Error(err),
					zap.String("instance_id", c.config.InstanceID.String()))
			}
		}
	}
}

func (c *client) sendEnqueuedMessagesLoop() {
	defer c.wg.Done()

	const batchInterval = 10 * time.Millisecond
	var messageBuffer []*websocketv1.Message
	var retryBuffer []*websocketv1.Message

	for {
		select {
		case <-c.stop:
			logger.Log.Debug("sendEnqueuedMessagesLoop stopping")
			if err := c.sendRemainingMessages(messageBuffer, retryBuffer); err != nil {
				logger.Log.Error("failed to send remaining messages during shutdown",
					zap.Error(err),
					zap.String("instance_id", c.config.InstanceID.String()))
			}
			return
		case msg, ok := <-c.messageQueue:
			if !ok {
				logger.Log.Debug("sendEnqueuedMessagesLoop stopping (queue closed)")
				if err := c.sendRemainingMessages(messageBuffer, retryBuffer); err != nil {
					logger.Log.Error("failed to send remaining messages on queue close",
						zap.Error(err),
						zap.String("instance_id", c.config.InstanceID.String()))
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
					if err := c.sendWithRetry(msg); err != nil {
						remainingMessages = append(remainingMessages, msg)
						logger.Log.Error("error sending message after retries",
							zap.Error(err),
							zap.String("message_id", msg.Id),
							zap.String("instance_id", c.config.InstanceID.String()))
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

func (c *client) sendWithRetry(msg *websocketv1.Message) error {
	var lastErr error
	for attempt := range maxMessageRetries {
		if attempt > 0 {
			delay := messageRetryDelay * time.Duration(1<<uint(attempt-1))
			logger.Log.Debug("retrying message send",
				zap.Int("attempt", attempt+1),
				zap.Duration("delay", delay),
				zap.String("message_id", msg.Id),
				zap.String("instance_id", c.config.InstanceID.String()))
			time.Sleep(delay)
		}

		if err := c.send(msg); err != nil {
			lastErr = err
			logger.Log.Warn("message send failed",
				zap.Error(err),
				zap.Int("attempt", attempt+1),
				zap.String("message_id", msg.Id),
				zap.String("instance_id", c.config.InstanceID.String()))
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to send message after %d attempts: %w", maxMessageRetries, lastErr)
}

func (c *client) sendRemainingMessages(messageBuffer, retryBuffer []*websocketv1.Message) error {
	if len(messageBuffer) == 0 && len(retryBuffer) == 0 {
		return nil
	}

	logger.Log.Info("sending remaining messages before shutdown",
		zap.Int("message_count", len(messageBuffer)+len(retryBuffer)),
		zap.String("instance_id", c.config.InstanceID.String()))

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Create a channel to signal completion
	done := make(chan error, 1)
	go func() {
		var err error
		for _, msg := range append(messageBuffer, retryBuffer...) {
			if err = c.sendWithRetry(msg); err != nil {
				break
			}
			time.Sleep(time.Millisecond)
		}
		done <- err
	}()

	// Wait for completion or timeout
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to send remaining messages: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout while sending remaining messages: %w", ctx.Err())
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
	var closeErr error
	c.shutdownOnce.Do(func() {
		logger.Log.Info("initiating client shutdown",
			zap.String("instance_id", c.config.InstanceID.String()))

		// 1. Signal all goroutines to stop
		close(c.stop)

		// 2. Close the message queue to prevent new messages
		close(c.messageQueue)

		// 3. Wait for goroutines to finish with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		shutdownDone := make(chan struct{})
		go func() {
			c.wg.Wait()
			close(shutdownDone)
		}()

		select {
		case <-shutdownDone:
			logger.Log.Debug("all goroutines stopped successfully")
		case <-shutdownCtx.Done():
			logger.Log.Warn("timeout waiting for goroutines to stop",
				zap.String("instance_id", c.config.InstanceID.String()))
		}

		// 4. Close the WebSocket connection
		c.connMu.Lock()
		if c.conn != nil {
			// Try to send a close message
			err := c.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Log.Warn("error sending close message",
					zap.Error(err),
					zap.String("instance_id", c.config.InstanceID.String()))
			}

			// Close the underlying connection
			if err := c.conn.Close(); err != nil {
				logger.Log.Error("error closing websocket connection",
					zap.Error(err),
					zap.String("instance_id", c.config.InstanceID.String()))
				closeErr = fmt.Errorf("failed to close websocket connection: %w", err)
			}
			c.conn = nil
		}
		c.connMu.Unlock()

		// 5. Clean up response channels
		c.respMu.Lock()
		for id, ch := range c.responses {
			close(ch)
			delete(c.responses, id)
		}
		c.responses = make(map[string]chan *websocketv1.Message)
		c.respMu.Unlock()

		// 6. Signal that the client is fully closed
		select {
		case c.done <- nil:
		default:
		}

		logger.Log.Info("client shutdown completed",
			zap.String("instance_id", c.config.InstanceID.String()))
	})

	return closeErr
}
