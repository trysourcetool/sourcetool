package ws

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
)

// Manager defines the interface for WebSocket connection management operations
// that are relevant to the domain or application layers.
type Manager interface {
	// SendToHost sends a message to a specific connected host instance via the pub/sub system.
	SendToHost(ctx context.Context, hostID uuid.UUID, msg *websocketv1.Message) error
	// SendToClient sends a message to a specific connected client session via the pub/sub system.
	SendToClient(ctx context.Context, sessionID uuid.UUID, msg *websocketv1.Message) error
	// Close gracefully shuts down the connection manager, closing all connections and stopping background processes.
	Close() error
	// PingConnectedHost pings a specific host instance to check if it is online.
	PingConnectedHost(hostID uuid.UUID) error
	// SetConnectedHost sets a specific host instance as connected.
	SetConnectedHost(hostInstance *hostinstance.HostInstance, apiKey *apikey.APIKey, conn *websocket.Conn)
	// DisconnectHost disconnects a specific host instance.
	DisconnectHost(hostID uuid.UUID)
	// SetConnectedClient sets a specific client session as connected.
	SetConnectedClient(session *session.Session, conn *websocket.Conn)
	// DisconnectClient disconnects a specific client session.
	DisconnectClient(sessionID uuid.UUID)
}
