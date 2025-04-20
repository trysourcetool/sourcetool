package ws

import (
	"net/http"

	"github.com/gorilla/websocket"

	wsSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	wsserver "github.com/trysourcetool/sourcetool/backend/internal/transport/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws/handlers"
)

func NewRouter(d *port.Dependencies) *wsserver.Router {
	middleware := NewMiddlewareEE(d.Repository)
	wsHandler := handlers.NewWebSocketHandler(
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		d.WSManager,
		wsSvc.NewServiceEE(d),
	)
	return wsserver.NewRouter(middleware, wsHandler)
}
