package ws

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/trysourcetool/sourcetool/backend/ee/ws"
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
		ws.NewServiceEE(d),
	)
	return wsserver.NewRouter(middleware, wsHandler)
}
