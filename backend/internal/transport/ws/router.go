package ws

import (
	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws/handlers"
)

type Router struct {
	middleware Middleware
	wsHandler  *handlers.WebSocketHandler
}

func NewRouter(middleware Middleware, wsHandler *handlers.WebSocketHandler) *Router {
	return &Router{
		middleware: middleware,
		wsHandler:  wsHandler,
	}
}

func (router *Router) Build() chi.Router {
	r := chi.NewRouter()
	r.With(router.middleware.Auth).HandleFunc("/", router.wsHandler.Handle)
	return r
}
