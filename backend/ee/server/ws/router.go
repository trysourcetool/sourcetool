package ws

import (
	"net/http"

	"github.com/gorilla/websocket"

	wsSvc "github.com/trysourcetool/sourcetool/backend/ee/ws/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
	wsserver "github.com/trysourcetool/sourcetool/backend/server/ws"
	"github.com/trysourcetool/sourcetool/backend/server/ws/handlers"
)

func NewRouter(d *infra.Dependency) *wsserver.Router {
	middleware := NewMiddlewareEE(d.Store)
	wsHandler := handlers.NewWebSocketHandler(
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		wsSvc.NewWebSocketServiceEE(d),
	)
	return wsserver.NewRouter(middleware, wsHandler)
}
