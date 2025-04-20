package wsmanager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
)

const (
	// Time allowed to read the next pong message from the client.
	clientPongWait = 1 * time.Minute

	// Time allowed to read the next pong message from the host.
	hostPongWait = 6 * time.Hour
)

// manager handles WebSocket connections for hosts and clients.
// It implements the port.WSManager interface.
type manager struct {
	connectedHosts   map[uuid.UUID]*connectedHost
	connectedClients map[uuid.UUID]*connectedClient
	hostsMutex       sync.RWMutex
	clientsMutex     sync.RWMutex
	pubsubClient     port.PubSub
	repo             port.Repository
	ctx              context.Context    // Context for managing goroutine lifecycle
	cancel           context.CancelFunc // Function to cancel the context
	wg               sync.WaitGroup     // WaitGroup to wait for goroutines to finish
}

// Compile-time check to ensure manager implements port.WSManager.
var _ port.WSManager = (*manager)(nil)

// NewManager creates and initializes a new WebSocket connection manager.
// It starts the background goroutines for handling pub/sub messages.
func NewManager(ctx context.Context, repo port.Repository, pubsubClient port.PubSub) port.WSManager {
	managerCtx, cancel := context.WithCancel(ctx) // Use parent context
	m := &manager{
		connectedHosts:   make(map[uuid.UUID]*connectedHost),
		connectedClients: make(map[uuid.UUID]*connectedClient),
		pubsubClient:     pubsubClient,
		repo:             repo,
		ctx:              managerCtx,
		cancel:           cancel,
	}

	m.wg.Add(2) // Add count for the two subscriber goroutines
	go m.subscribeToHostMessages()
	go m.subscribeToClientMessages()

	return m
}

// SendToHost publishes a message destined for a specific host to the pub/sub system.
func (m *manager) SendToHost(ctx context.Context, hostInstanceID uuid.UUID, msg *websocketv1.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message for host %s: %w", hostInstanceID, err)
	}

	logger.Logger.Sugar().Debugf("Publishing message %s to host_messages for host %s", msg.Id, hostInstanceID)
	return m.pubsubClient.Publish(ctx, "host_messages", hostInstanceID.String(), data)
}

// SendToClient publishes a message destined for a specific client to the pub/sub system.
func (m *manager) SendToClient(ctx context.Context, sessionID uuid.UUID, msg *websocketv1.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message for client %s: %w", sessionID, err)
	}

	logger.Logger.Sugar().Debugf("Publishing message %s to client_messages for client %s", msg.Id, sessionID)
	return m.pubsubClient.Publish(ctx, "client_messages", sessionID.String(), data)
}

// SetConnectedHost registers a new host connection with the manager.
// It handles potential existing connections and starts the ping loop.
func (m *manager) SetConnectedHost(hostInstance *hostinstance.HostInstance, apiKey *apikey.APIKey, conn *websocket.Conn) {
	// Disconnect any existing connection for the same host ID first.
	m.DisconnectHost(hostInstance.ID)

	logger.Logger.Sugar().Infof("Registering new connection for host: %s", hostInstance.ID)

	conn.SetPongHandler(func(string) error {
		logger.Logger.Sugar().Debugf("Received pong from host %s", hostInstance.ID)
		return conn.SetReadDeadline(time.Now().Add(hostPongWait))
	})

	host := &connectedHost{
		hostInstance: hostInstance,
		apiKey:       apiKey,
		conn:         conn,
		done:         make(chan struct{}),
	}

	m.hostsMutex.Lock()
	m.connectedHosts[hostInstance.ID] = host
	m.hostsMutex.Unlock()

	m.wg.Add(1) // Add to WaitGroup for the ping loop
	go func() {
		defer m.wg.Done()
		m.startHostPingLoop(host)
	}()
}

// DisconnectHost removes a host connection from the manager and closes the connection.
func (m *manager) DisconnectHost(hostInstanceID uuid.UUID) {
	m.hostsMutex.Lock()
	host, ok := m.connectedHosts[hostInstanceID]
	if ok {
		delete(m.connectedHosts, hostInstanceID)
	}
	m.hostsMutex.Unlock()

	if ok {
		close(host.done) // Stop ping loop
		logger.Logger.Sugar().Infof("Disconnecting host: %s", hostInstanceID)

		// Explicitly close the WebSocket connection
		if err := host.conn.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Failed to close host WebSocket connection for %s: %v", hostInstanceID, err)
		} else {
			logger.Logger.Sugar().Debugf("Closed host WebSocket connection for %s", hostInstanceID)
		}
	} else {
		logger.Logger.Sugar().Debugf("Attempted to disconnect host %s, but it was not found.", hostInstanceID)
	}
}

// SetConnectedClient registers a new client connection with the manager.
// It handles potential existing connections and starts the ping loop.
func (m *manager) SetConnectedClient(session *session.Session, conn *websocket.Conn) {
	// Disconnect any existing connection for the same session ID first.
	m.DisconnectClient(session.ID)

	logger.Logger.Sugar().Infof("Registering new connection for client session: %s", session.ID)

	conn.SetPongHandler(func(string) error {
		logger.Logger.Sugar().Debugf("Received pong from client %s", session.ID)
		return conn.SetReadDeadline(time.Now().Add(clientPongWait))
	})

	client := &connectedClient{
		session: session,
		conn:    conn,
		done:    make(chan struct{}),
	}

	m.clientsMutex.Lock()
	m.connectedClients[session.ID] = client
	m.clientsMutex.Unlock()

	m.wg.Add(1) // Add to WaitGroup for the ping loop
	go func() {
		defer m.wg.Done()
		m.startClientPingLoop(client)
	}()
}

// DisconnectClient removes a client connection from the manager and closes the connection.
func (m *manager) DisconnectClient(sessionID uuid.UUID) {
	m.clientsMutex.Lock()
	client, ok := m.connectedClients[sessionID]
	if ok {
		delete(m.connectedClients, sessionID)
	}
	m.clientsMutex.Unlock()

	if ok {
		close(client.done) // Stop ping loop
		logger.Logger.Sugar().Infof("Disconnecting client: %s", sessionID)

		// Explicitly close the WebSocket connection
		if err := client.conn.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Failed to close client WebSocket connection for %s: %v", sessionID, err)
		} else {
			logger.Logger.Sugar().Debugf("Closed client WebSocket connection for %s", sessionID)
		}
	} else {
		logger.Logger.Sugar().Debugf("Attempted to disconnect client %s, but it was not found.", sessionID)
	}
}

// Close gracefully shuts down the connection manager.
// It stops all background goroutines (ping loops, subscribers) and closes connections.
func (m *manager) Close() error {
	logger.Logger.Sugar().Info("Closing WebSocket connection manager...")

	logger.Logger.Sugar().Info("Canceling connection manager context...")
	m.cancel() // This signals subscribers to stop

	logger.Logger.Sugar().Info("Signaling ping loops to stop...")
	m.hostsMutex.Lock()
	for id, host := range m.connectedHosts {
		logger.Logger.Sugar().Debugf("Closing done channel for host %s", id)
		close(host.done)
		if err := host.conn.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Failed to close host WebSocket connection for %s: %v", id, err)
		} else {
			logger.Logger.Sugar().Debugf("Closed host WebSocket connection for %s", id)
		}
	}
	m.hostsMutex.Unlock()

	m.clientsMutex.Lock()
	for id, client := range m.connectedClients {
		logger.Logger.Sugar().Debugf("Closing done channel for client %s", id)
		close(client.done)
		if err := client.conn.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Failed to close client WebSocket connection for %s: %v", id, err)
		} else {
			logger.Logger.Sugar().Debugf("Closed client WebSocket connection for %s", id)
		}
	}
	m.clientsMutex.Unlock()

	logger.Logger.Sugar().Info("Waiting for background goroutines to stop...")
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Logger.Sugar().Info("All background goroutines stopped.")
	case <-time.After(30 * time.Second): // Add a timeout
		logger.Logger.Sugar().Error("Timeout waiting for background goroutines to stop.")
		return errors.New("timeout closing connection manager goroutines")
	}

	m.hostsMutex.Lock()
	m.connectedHosts = make(map[uuid.UUID]*connectedHost)
	m.hostsMutex.Unlock()

	m.clientsMutex.Lock()
	m.connectedClients = make(map[uuid.UUID]*connectedClient)
	m.clientsMutex.Unlock()

	logger.Logger.Sugar().Info("WebSocket connection manager closed.")
	return nil
}

// PingConnectedHost is an internal method (not part of port.WSManager) for sending a direct ping.
// It's kept separate as direct pinging might be an infra-specific detail.
// Renamed to PingConnectedHost to avoid conflict with interface methods if any.
func (m *manager) PingConnectedHost(hostInstanceID uuid.UUID) error {
	m.hostsMutex.RLock()
	host, ok := m.connectedHosts[hostInstanceID]
	m.hostsMutex.RUnlock()

	if ok {
		return m.pingConnection(host.conn)
	}
	return fmt.Errorf("host connection %s not found for ping", hostInstanceID)
}
