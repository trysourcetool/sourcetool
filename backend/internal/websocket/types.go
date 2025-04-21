package websocket

import (
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

// connectedHost represents a connected host instance with its associated data.
type connectedHost struct {
	hostInstance *core.HostInstance
	apiKey       *core.APIKey
	conn         *websocket.Conn
	done         chan struct{} // Channel to signal termination for the host's goroutines (e.g., ping loop)
}

// connectedClient represents a connected client session with its associated data.
type connectedClient struct {
	session *core.Session
	conn    *websocket.Conn
	done    chan struct{} // Channel to signal termination for the client's goroutines (e.g., ping loop)
}
