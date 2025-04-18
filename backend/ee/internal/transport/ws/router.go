package ws

import (
	"net/http"

	"github.com/gorilla/websocket"

	wsSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	wsserver "github.com/trysourcetool/sourcetool/backend/internal/transport/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws/handlers"
)

func NewRouter(d *infra.Dependency) *wsserver.Router {
	middleware := NewMiddlewareEE(d.Repository)
	wsHandler := handlers.NewWebSocketHandler(
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		wsSvc.NewServiceEE(d),
	)
	return wsserver.NewRouter(middleware, wsHandler)
}
