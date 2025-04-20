package wsmanager

import (
	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
)

// connectedHost represents a connected host instance with its associated data.
type connectedHost struct {
	hostInstance *hostinstance.HostInstance
	apiKey       *apikey.APIKey
	conn         *websocket.Conn
	done         chan struct{} // Channel to signal termination for the host's goroutines (e.g., ping loop)
}

// connectedClient represents a connected client session with its associated data.
type connectedClient struct {
	session *session.Session
	conn    *websocket.Conn
	done    chan struct{} // Channel to signal termination for the client's goroutines (e.g., ping loop)
}
