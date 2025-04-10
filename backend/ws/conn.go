package ws

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
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
	connectedHosts   *sync.Map // hostInstanceID -> connectedHost
	connectedClients *sync.Map // sessionID -> connectedClient
	redisClient      *redisClient
	store            infra.Store
	ctx              context.Context    // Context for managing goroutine lifecycle
	cancel           context.CancelFunc // Function to cancel the context
	wg               sync.WaitGroup     // WaitGroup to wait for goroutines to finish
}

func newConnManager(redisClient *redisClient, store infra.Store) *connManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &connManager{
		connectedHosts:   new(sync.Map),
		connectedClients: new(sync.Map),
		redisClient:      redisClient,
		store:            store,
		ctx:              ctx,
		cancel:           cancel,
	}
}

func (c *connManager) pingConnection(conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeWait)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
				logger.Logger.Sugar().Errorf("Failed to ping client %s: %v", client.session.ID, err)

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
	defer c.wg.Done() // Signal WaitGroup when this goroutine exits

	// Use the manager's context for the subscription
	ch, err := c.redisClient.Subscribe(c.ctx, "host_messages")
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to subscribe to host messages: %v", err)
		// Consider more robust error handling here (retry, health check failure, etc.)
		return
	}
	logger.Logger.Sugar().Info("Subscribed to host messages")

	for {
		select {
		case <-c.ctx.Done(): // Exit loop if manager context is canceled
			logger.Logger.Sugar().Info("Host message subscriber stopping due to context cancellation.")
			// The pubsub channel might be closed implicitly by redis client on context cancel/close,
			// but explicitly returning ensures termination.
			return
		case msg, ok := <-ch:
			if !ok {
				logger.Logger.Sugar().Info("Host message channel closed.")
				return // Channel closed, exit goroutine
			}

			// Existing message processing logic...
			var redisMsg redisv1.RedisMessage
			if err := proto.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to unmarshal redis message: %v", err)
				continue
			}

			hostInstanceID, err := uuid.FromString(redisMsg.Id)
			if err != nil {
				logger.Logger.Sugar().Errorf("Invalid host instance ID: %v", err)
				continue
			}

			var protoMsg websocketv1.Message
			if err := proto.Unmarshal(redisMsg.Payload, &protoMsg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to unmarshal protobuf message: %v", err)
				continue
			}

			logger.Logger.Sugar().Debugf("Received message: %s", &protoMsg)

			conn, ok := c.connectedHosts.Load(hostInstanceID)
			if !ok {
				// Host might have disconnected, log as debug or info?
				logger.Logger.Sugar().Debugf("Host not found for message: %s", hostInstanceID)
				continue
			}
			host, ok := conn.(*connectedHost)
			if !ok {
				logger.Logger.Sugar().Errorf("Invalid connection type stored for host: %s", hostInstanceID)
				c.connectedHosts.Delete(hostInstanceID) // Clean up invalid entry
				continue
			}

			logger.Logger.Sugar().Debugf("Sending message to host %s: %s", host.hostInstance.ID, protoMsg.Id)

			data, err := proto.Marshal(&protoMsg)
			if err != nil {
				logger.Logger.Sugar().Errorf("Failed to marshal message for host %s: %v", host.hostInstance.ID, err)
				continue
			}

			// Consider adding a write deadline
			if err := host.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				logger.Logger.Sugar().Errorf("Failed to send message to host %s: %v. Disconnecting host.", host.hostInstance.ID, err)
				// Consider attempting to disconnect the host here if writing fails repeatedly
				// c.DisconnectHost(hostInstanceID)
			}
		}
	}
}

func (c *connManager) subscribeToClientMessages() {
	defer c.wg.Done() // Signal WaitGroup when this goroutine exits

	// Use the manager's context for the subscription
	ch, err := c.redisClient.Subscribe(c.ctx, "client_messages")
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to subscribe to client messages: %v", err)
		// Consider more robust error handling here
		return
	}
	logger.Logger.Sugar().Info("Subscribed to client messages")

	for {
		select {
		case <-c.ctx.Done(): // Exit loop if manager context is canceled
			logger.Logger.Sugar().Info("Client message subscriber stopping due to context cancellation.")
			return
		case msg, ok := <-ch:
			if !ok {
				logger.Logger.Sugar().Info("Client message channel closed.")
				return // Channel closed, exit goroutine
			}

			// Existing message processing logic...
			var redisMsg redisv1.RedisMessage
			if err := proto.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to unmarshal redis message: %v", err)
				continue
			}

			sessionID, err := uuid.FromString(redisMsg.Id)
			if err != nil {
				logger.Logger.Sugar().Errorf("Invalid session ID: %v", err)
				continue
			}

			var protoMsg websocketv1.Message
			if err := proto.Unmarshal(redisMsg.Payload, &protoMsg); err != nil {
				logger.Logger.Sugar().Errorf("Failed to unmarshal protobuf message: %v", err)
				continue
			}

			logger.Logger.Sugar().Debugf("Received message: %s", &protoMsg)

			conn, ok := c.connectedClients.Load(sessionID)
			if !ok {
				// Client might have disconnected
				logger.Logger.Sugar().Debugf("Client not found for message: %s", sessionID)
				continue
			}
			client, ok := conn.(*connectedClient)
			if !ok {
				logger.Logger.Sugar().Errorf("Invalid connection type stored for client: %s", sessionID)
				c.connectedClients.Delete(sessionID) // Clean up invalid entry
				continue
			}

			logger.Logger.Sugar().Debugf("Sending message to client %s: %s", client.session.ID, protoMsg.Id)

			data, err := proto.Marshal(&protoMsg)
			if err != nil {
				logger.Logger.Sugar().Errorf("Failed to marshal message for client %s: %v", client.session.ID, err)
				continue
			}

			// Consider adding a write deadline
			if err := client.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				logger.Logger.Sugar().Errorf("Failed to send message to client %s: %v. Disconnecting client.", client.session.ID, err)
				// Consider attempting to disconnect the client here
				// c.DisconnectClient(sessionID)
			}
		}
	}
}

func (c *connManager) PingHost(hostInstanceID uuid.UUID) error {
	conn, ok := c.connectedHosts.Load(hostInstanceID)
	if !ok {
		return errors.New("connection not found")
	}
	connectedHost, ok := conn.(*connectedHost)
	if !ok {
		return errors.New("invalid connection")
	}
	return c.pingConnection(connectedHost.conn)
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

	c.connectedHosts.Store(hostInstance.ID, host)

	go c.startHostPingLoop(host)
}

func (c *connManager) DisconnectHost(hostInstanceID uuid.UUID) {
	if conn, ok := c.connectedHosts.Load(hostInstanceID); ok {
		if host, ok := conn.(*connectedHost); ok {
			close(host.done) // Stop ping loop
			logger.Logger.Sugar().Debug("Stopped ping host")
		}
	}
	c.connectedHosts.Delete(hostInstanceID)
}

func (c *connManager) SetConnectedClient(session *model.Session, conn *websocket.Conn) {
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(clientPongWait))
	})

	client := &connectedClient{
		session: session,
		conn:    conn,
		done:    make(chan struct{}),
	}

	c.connectedClients.Store(session.ID, client)

	go c.startClientPingLoop(client)
}

func (c *connManager) DisconnectClient(sessionID uuid.UUID) {
	if conn, ok := c.connectedClients.Load(sessionID); ok {
		if client, ok := conn.(*connectedClient); ok {
			close(client.done) // Stop ping loop
			logger.Logger.Sugar().Debug("Stopped ping client")
		}
	}
	c.connectedClients.Delete(sessionID)
}

func (c *connManager) Close() {
	logger.Logger.Sugar().Info("Closing connection manager...")

	// 1. Stop accepting new connections (implicitly done by server shutdown)
	// 2. Stop ping loops for existing connections
	c.connectedHosts.Range(func(key, value interface{}) bool {
		if host, ok := value.(*connectedHost); ok {
			close(host.done) // Signal ping loop to stop
		}
		// Keep host in map for now, might receive final messages from Redis
		return true
	})
	c.connectedClients.Range(func(key, value interface{}) bool {
		if client, ok := value.(*connectedClient); ok {
			close(client.done) // Signal ping loop to stop
		}
		// Keep client in map for now
		return true
	})

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

	// 6. Optionally close actual WebSocket connections (gorilla/websocket handles this often)
	//    Or just clear the maps now.
	c.connectedHosts = new(sync.Map)
	c.connectedClients = new(sync.Map)

	logger.Logger.Sugar().Info("Connection manager closed.")
}
