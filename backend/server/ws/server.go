package ws

import (
	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/server/ws/handlers"
)

type ServerCE struct {
	middleware MiddlewareCE
	wsHandler  handlers.WebSocketHandlerCE
}

func NewServerCE(middleware MiddlewareCE, wsHandler handlers.WebSocketHandlerCE) *ServerCE {
	return &ServerCE{
		middleware: middleware,
		wsHandler:  wsHandler,
	}
}

func (s *ServerCE) Router() chi.Router {
	r := chi.NewRouter()
	r.With(s.middleware.Auth).HandleFunc("/", s.wsHandler.Handle)
	return r
}
