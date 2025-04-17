package ws

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
	redisv1 "github.com/trysourcetool/sourcetool/backend/pb/go/redis/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/pb/go/websocket/v1"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Send pings to peer with this period.
	pingPeriod = 30 * time.Second

	// Time allowed to read the next pong message from the client.
	clientPongWait = 1 * time.Minute

	// Time allowed to read the next pong message from the host.
	hostPongWait = 6 * time.Hour

	// Maximum number of reconnection attempts.
	maxReconnectAttempts = 5

	// Base delay for exponential backoff (in milliseconds).
	baseReconnectDelay = 100
)

var (
	connManagerInstance *connManager
	once                sync.Once
)

func GetConnManager() *connManager {
	return connManagerInstance
}

type connectedHost struct {
	hostInstance *model.HostInstance
	apiKey       *model.APIKey
	conn         *websocket.Conn
	done         chan struct{}
}

type connectedClient struct {
	session *model.Session
	conn    *websocket.Conn
	done    chan struct{}
}

type connManager struct {
	connectedHosts   map[uuid.UUID]*connectedHost
	connectedClients map[uuid.UUID]*connectedClient
	hostsMutex       sync.RWMutex
	clientsMutex     sync.RWMutex
	redisClient      *redisClient
	store            infra.Store
	ctx              context.Context    // Context for managing goroutine lifecycle
	cancel           context.CancelFunc // Function to cancel the context
	wg               sync.WaitGroup     // WaitGroup to wait for goroutines to finish
}

func newConnManager(redisClient *redisClient, store infra.Store) *connManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &connManager{
		connectedHosts:   make(map[uuid.UUID]*connectedHost),
		connectedClients: make(map[uuid.UUID]*connectedClient),
		redisClient:      redisClient,
		store:            store,
		ctx:              ctx,
		cancel:           cancel,
	}
}

func (c *connManager) pingConnection(conn *websocket.Conn) error {
	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Write ping message directly
	if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to write ping message: %w", err)
	}

	return nil
}

func (c *connManager) startHostPingLoop(host *connectedHost) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-host.done:
			return
		case <-ticker.C:
			if err := c.pingConnection(host.conn); err != nil {
				// Consider retrying the connection a few times before disconnecting the host immediately.

				logger.Logger.Sugar().Errorf("Failed to ping host %s: %v", host.hostInstance.ID, err)

				host.hostInstance.Status = model.HostInstanceStatusUnreachable

				if err := c.store.HostInstance().Update(context.Background(), host.hostInstance); err != nil {
					logger.Logger.Sugar().Errorf("Failed to update host status: %v", err)
				}

				c.DisconnectHost(host.hostInstance.ID)
				return
			}

			logger.Logger.Sugar().Debug("Successfully pinged host")
		}
	}
}

func (c *connManager) startClientPingLoop(client *connectedClient) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-client.done:
			return
		case <-ticker.C:
			if err := c.pingConnection(client.conn); err != nil {
				// Consider retrying the connection a few times before disconnecting the client immediately.

				logger.Logger.Sugar().Errorf("Failed to ping client %s: %v", client.session.ID, err)

				if err := c.store.Session().Delete(context.Background(), client.session); err != nil {
					logger.Logger.Sugar().Errorf("Failed to delete session %s: %v", client.session.ID, err)
				}

				c.DisconnectClient(client.session.ID)
				return
			}

			logger.Logger.Sugar().Debug("Successfully pinged client")
		}
	}
}

func InitWebSocketConns(ctx context.Context, store infra.Store) error {
	var initErr error
	once.Do(func() {
		redisClient, err := newRedisClient()
		if err != nil {
			initErr = err
			return
		}
		connManagerInstance = newConnManager(redisClient, store)

		connManagerInstance.wg.Add(2) // Add count for the two subscriber goroutines
		go connManagerInstance.subscribeToHostMessages()
		go connManagerInstance.subscribeToClientMessages()
	})

	return initErr
}

func (c *connManager) subscribeToHostMessages() {
	defer c.wg.Done()

	for attempt := range maxReconnectAttempts {
		if err := c.subscribeToHostMessagesWithRetry(); err != nil {
			if attempt == maxReconnectAttempts-1 {
				logger.Logger.Sugar().Errorf("Failed to subscribe to host messages after %d attempts: %v", maxReconnectAttempts, err)
				return
			}
			// Exponential backoff
			delay := time.Duration(baseReconnectDelay*(1<<attempt)) * time.Millisecond
			logger.Logger.Sugar().Warnf("Retrying host message subscription in %v (attempt %d/%d)", delay, attempt+1, maxReconnectAttempts)
			time.Sleep(delay)
			continue
		}
		return
	}
}

func (c *connManager) subscribeToHostMessagesWithRetry() error {
	ch, err := c.redisClient.Subscribe(c.ctx, "host_messages")
	if err != nil {
		return fmt.Errorf("failed to subscribe to host messages: %w", err)
	}
	logger.Logger.Sugar().Info("Subscribed to host messages")

	for {
		select {
		case <-c.ctx.Done():
			logger.Logger.Sugar().Info("Host message subscriber stopping due to context cancellation.")
			return nil
		case msg, ok := <-ch:
			if !ok {
				return fmt.Errorf("host message channel closed unexpectedly")
			}

			if err := c.processHostMessage(msg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to process host message: %v", err)
				// Continue processing other messages even if one fails
				continue
			}
		}
	}
}

func (c *connManager) processHostMessage(msg *redis.Message) error {
	var redisMsg redisv1.RedisMessage
	if err := proto.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
		return fmt.Errorf("failed to unmarshal redis message: %w", err)
	}

	hostInstanceID, err := uuid.FromString(redisMsg.Id)
	if err != nil {
		return fmt.Errorf("invalid host instance ID: %w", err)
	}

	var protoMsg websocketv1.Message
	if err := proto.Unmarshal(redisMsg.Payload, &protoMsg); err != nil {
		return fmt.Errorf("failed to unmarshal protobuf message: %w", err)
	}

	logger.Logger.Sugar().Debugf("Received message: %s", &protoMsg)

	c.hostsMutex.RLock()
	host, ok := c.connectedHosts[hostInstanceID]
	c.hostsMutex.RUnlock()

	if !ok {
		logger.Logger.Sugar().Debugf("Host not found for message: %s", hostInstanceID)
		return nil
	}

	logger.Logger.Sugar().Debugf("Sending message to host %s: %s", host.hostInstance.ID, protoMsg.Id)

	data, err := proto.Marshal(&protoMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message for host %s: %w", host.hostInstance.ID, err)
	}

	if err := host.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to set write deadline for host %s: %w", host.hostInstance.ID, err)
	}

	if err := host.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		c.DisconnectHost(hostInstanceID)
		return fmt.Errorf("failed to send message to host %s: %w", host.hostInstance.ID, err)
	}

	return nil
}

func (c *connManager) subscribeToClientMessages() {
	defer c.wg.Done()

	for attempt := range maxReconnectAttempts {
		if err := c.subscribeToClientMessagesWithRetry(); err != nil {
			if attempt == maxReconnectAttempts-1 {
				logger.Logger.Sugar().Errorf("Failed to subscribe to client messages after %d attempts: %v", maxReconnectAttempts, err)
				return
			}
			// Exponential backoff
			delay := time.Duration(baseReconnectDelay*(1<<attempt)) * time.Millisecond
			logger.Logger.Sugar().Warnf("Retrying client message subscription in %v (attempt %d/%d)", delay, attempt+1, maxReconnectAttempts)
			time.Sleep(delay)
			continue
		}
		return
	}
}

func (c *connManager) subscribeToClientMessagesWithRetry() error {
	ch, err := c.redisClient.Subscribe(c.ctx, "client_messages")
	if err != nil {
		return fmt.Errorf("failed to subscribe to client messages: %w", err)
	}
	logger.Logger.Sugar().Info("Subscribed to client messages")

	for {
		select {
		case <-c.ctx.Done():
			logger.Logger.Sugar().Info("Client message subscriber stopping due to context cancellation.")
			return nil
		case msg, ok := <-ch:
			if !ok {
				return fmt.Errorf("client message channel closed unexpectedly")
			}

			if err := c.processClientMessage(msg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to process client message: %v", err)
				// Continue processing other messages even if one fails
				continue
			}
		}
	}
}

func (c *connManager) processClientMessage(msg *redis.Message) error {
	var redisMsg redisv1.RedisMessage
	if err := proto.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
		return fmt.Errorf("failed to unmarshal redis message: %w", err)
	}

	sessionID, err := uuid.FromString(redisMsg.Id)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	var protoMsg websocketv1.Message
	if err := proto.Unmarshal(redisMsg.Payload, &protoMsg); err != nil {
		return fmt.Errorf("failed to unmarshal protobuf message: %w", err)
	}

	logger.Logger.Sugar().Debugf("Received message: %s", &protoMsg)

	c.clientsMutex.RLock()
	client, ok := c.connectedClients[sessionID]
	c.clientsMutex.RUnlock()

	if !ok {
		logger.Logger.Sugar().Debugf("Client not found for message: %s", sessionID)
		return nil
	}

	logger.Logger.Sugar().Debugf("Sending message to client %s: %s", client.session.ID, protoMsg.Id)

	data, err := proto.Marshal(&protoMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message for client %s: %w", client.session.ID, err)
	}

	if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to set write deadline for client %s: %w", client.session.ID, err)
	}

	if err := client.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		c.DisconnectClient(sessionID)
		return fmt.Errorf("failed to send message to client %s: %w", client.session.ID, err)
	}

	return nil
}

func (c *connManager) PingHost(hostInstanceID uuid.UUID) error {
	c.hostsMutex.RLock()
	defer c.hostsMutex.RUnlock()

	if host, ok := c.connectedHosts[hostInstanceID]; ok {
		return c.pingConnection(host.conn)
	}
	return errors.New("connection not found")
}

func (c *connManager) SendToHost(ctx context.Context, hostInstanceID uuid.UUID, msg *websocketv1.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	return c.redisClient.Publish(ctx, "host_messages", hostInstanceID.String(), data)
}

func (c *connManager) SendToClient(ctx context.Context, sessionID uuid.UUID, msg *websocketv1.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	return c.redisClient.Publish(ctx, "client_messages", sessionID.String(), data)
}

func (c *connManager) SetConnectedHost(hostInstance *model.HostInstance, apiKey *model.APIKey, conn *websocket.Conn) {
	// If a connection with the same hostInstance.ID already exists, disconnect it first to avoid duplicate connections.
	c.hostsMutex.Lock()
	if _, exists := c.connectedHosts[hostInstance.ID]; exists {
		c.hostsMutex.Unlock()
		c.DisconnectHost(hostInstance.ID)
		c.hostsMutex.Lock()
	}
	c.hostsMutex.Unlock()
	logger.Logger.Sugar().Debugf("Connected host: %s", hostInstance.ID)

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(hostPongWait))
	})

	host := &connectedHost{
		hostInstance: hostInstance,
		apiKey:       apiKey,
		conn:         conn,
		done:         make(chan struct{}),
	}

	c.hostsMutex.Lock()
	c.connectedHosts[hostInstance.ID] = host
	c.hostsMutex.Unlock()

	go c.startHostPingLoop(host)
}

func (c *connManager) DisconnectHost(hostInstanceID uuid.UUID) {
	c.hostsMutex.Lock()
	defer c.hostsMutex.Unlock()

	if host, ok := c.connectedHosts[hostInstanceID]; ok {
		close(host.done) // Stop ping loop
		logger.Logger.Sugar().Debug("Stopped ping host")

		// Explicitly close the WebSocket connection
		if err := host.conn.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Failed to close host WebSocket connection: %v", err)
		} else {
			logger.Logger.Sugar().Debugf("Closed host WebSocket connection for host %s", hostInstanceID)
		}
		delete(c.connectedHosts, hostInstanceID)
	}
}

func (c *connManager) SetConnectedClient(session *model.Session, conn *websocket.Conn) {
	// If a connection with the same session.ID already exists, disconnect it first to avoid duplicate connections.
	c.clientsMutex.Lock()
	if _, exists := c.connectedClients[session.ID]; exists {
		c.clientsMutex.Unlock() // Unlock before calling DisconnectClient (which will re-lock internally)
		c.DisconnectClient(session.ID)
		c.clientsMutex.Lock() // Re-lock after disconnect
	}
	c.clientsMutex.Unlock()
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(clientPongWait))
	})

	client := &connectedClient{
		session: session,
		conn:    conn,
		done:    make(chan struct{}),
	}

	c.clientsMutex.Lock()
	c.connectedClients[session.ID] = client
	c.clientsMutex.Unlock()

	go c.startClientPingLoop(client)
}

func (c *connManager) DisconnectClient(sessionID uuid.UUID) {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	if client, ok := c.connectedClients[sessionID]; ok {
		close(client.done) // Stop ping loop
		logger.Logger.Sugar().Debug("Stopped ping client")

		// Explicitly close the WebSocket connection
		if err := client.conn.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Failed to close client WebSocket connection: %v", err)
		} else {
			logger.Logger.Sugar().Debugf("Closed client WebSocket connection for session %s", sessionID)
		}
		delete(c.connectedClients, sessionID)
	}
}

func (c *connManager) Close() {
	logger.Logger.Sugar().Info("Closing connection manager...")

	// 1. Stop accepting new connections (implicitly done by server shutdown)
	// 2. Stop ping loops for existing connections
	c.hostsMutex.Lock()
	for _, host := range c.connectedHosts {
		close(host.done) // Signal ping loop to stop
	}
	c.hostsMutex.Unlock()

	c.clientsMutex.Lock()
	for _, client := range c.connectedClients {
		close(client.done) // Signal ping loop to stop
	}
	c.clientsMutex.Unlock()

	// 3. Signal Redis subscriber goroutines to stop
	logger.Logger.Sugar().Info("Canceling connection manager context...")
	c.cancel()

	// 4. Wait for subscriber goroutines to finish
	logger.Logger.Sugar().Info("Waiting for subscriber goroutines to stop...")
	c.wg.Wait()
	logger.Logger.Sugar().Info("Subscriber goroutines stopped.")

	// 5. Close the Redis client connection
	logger.Logger.Sugar().Info("Closing Redis client connection...")
	if err := c.redisClient.Close(); err != nil {
		logger.Logger.Sugar().Errorf("Failed to close Redis client: %v", err)
	}

	// 6. Clear the maps
	c.hostsMutex.Lock()
	c.connectedHosts = make(map[uuid.UUID]*connectedHost)
	c.hostsMutex.Unlock()

	c.clientsMutex.Lock()
	c.connectedClients = make(map[uuid.UUID]*connectedClient)
	c.clientsMutex.Unlock()

	logger.Logger.Sugar().Info("Connection manager closed.")
}
