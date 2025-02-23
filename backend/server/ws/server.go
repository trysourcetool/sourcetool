package ws

import (
	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/server/ws/handlers"
)

type Server struct {
	middleware Middleware
	wsHandler  handlers.WebSocketHandler
}

func NewServer(middleware Middleware, wsHandler handlers.WebSocketHandler) *Server {
	return &Server{
		middleware: middleware,
		wsHandler:  wsHandler,
	}
}

func (s *Server) Router() chi.Router {
	r := chi.NewRouter()
	r.With(s.middleware.Auth).HandleFunc("/", s.wsHandler.Handle)
	return r
}
