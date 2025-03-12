package ws

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	redisv1 "github.com/trysourcetool/sourcetool/proto/go/redis/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
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
}

func newConnManager(redisClient *redisClient, store infra.Store) *connManager {
	return &connManager{
		connectedHosts:   new(sync.Map),
		connectedClients: new(sync.Map),
		redisClient:      redisClient,
		store:            store,
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

		go connManagerInstance.subscribeToHostMessages(ctx)
		go connManagerInstance.subscribeToClientMessages(ctx)
	})

	return initErr
}

func (c *connManager) subscribeToHostMessages(ctx context.Context) {
	ch, err := c.redisClient.Subscribe(ctx, "host_messages")
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to subscribe to host messages: %v", err)
		return
	}

	for msg := range ch {
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
			logger.Logger.Sugar().Errorf("Host not found: %s", hostInstanceID)
			continue
		}
		host, ok := conn.(*connectedHost)
		if !ok {
			logger.Logger.Sugar().Errorf("Invalid connection: %s", hostInstanceID)
			continue
		}

		logger.Logger.Sugar().Debugf("Sending message to host %s: %s", host.hostInstance.ID, protoMsg.Id)

		data, err := proto.Marshal(&protoMsg)
		if err != nil {
			logger.Logger.Sugar().Errorf("Failed to marshal message for host %s: %v", host.hostInstance.ID, err)
			continue
		}

		if err := host.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			logger.Logger.Sugar().Errorf("Failed to send message to host %s: %v", host.hostInstance.ID, err)
		}
	}
}

func (c *connManager) subscribeToClientMessages(ctx context.Context) {
	ch, err := c.redisClient.Subscribe(ctx, "client_messages")
	if err != nil {
		logger.Logger.Sugar().Errorf("Failed to subscribe to client messages: %v", err)
		return
	}

	for msg := range ch {
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
			logger.Logger.Sugar().Errorf("Client not found: %s", sessionID)
			continue
		}
		client, ok := conn.(*connectedClient)
		if !ok {
			logger.Logger.Sugar().Errorf("Invalid connection: %s", sessionID)
			continue
		}

		logger.Logger.Sugar().Debugf("Sending message to client %s: %s", client.session.ID, protoMsg.Id)

		data, err := proto.Marshal(&protoMsg)
		if err != nil {
			logger.Logger.Sugar().Errorf("Failed to marshal message for client %s: %v", client.session.ID, err)
			continue
		}

		if err := client.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			logger.Logger.Sugar().Errorf("Failed to send message to client %s: %v", client.session.ID, err)
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
	c.connectedHosts.Range(func(key, value interface{}) bool {
		if host, ok := value.(*connectedHost); ok {
			close(host.done)
		}
		return true
	})

	c.connectedClients.Range(func(key, value interface{}) bool {
		if client, ok := value.(*connectedClient); ok {
			close(client.done)
		}
		return true
	})
}
