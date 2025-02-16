package server

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/page"
	httpserver "github.com/trysourcetool/sourcetool/backend/server/http"
	httphandlers "github.com/trysourcetool/sourcetool/backend/server/http/handlers"
	"github.com/trysourcetool/sourcetool/backend/server/ws"
	wshandlers "github.com/trysourcetool/sourcetool/backend/server/ws/handlers"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type Server struct {
	wsServer   *ws.ServerCE
	httpServer *httpserver.ServerCE
}

func New(d *infra.Dependency) *Server {
	httpMiddle := httpserver.NewMiddlewareCE(d.Store)
	apiKeyHandler := httphandlers.NewAPIKeyHandlerCE(apikey.NewServiceCE(d))
	environmentHandler := httphandlers.NewEnvironmentHandlerCE(environment.NewServiceCE(d))
	groupHandler := httphandlers.NewGroupHandlerCE(group.NewServiceCE(d))
	hostInstanceHandler := httphandlers.NewHostInstanceHandlerCE(hostinstance.NewServiceCE(d))
	organizationHandler := httphandlers.NewOrganizationHandlerCE(organization.NewServiceCE(d))
	pageHandler := httphandlers.NewPageHandlerCE(page.NewServiceCE(d))
	userHandler := httphandlers.NewUserHandlerCE(user.NewServiceCE(d))

	wsMiddle := ws.NewMiddlewareCE(d.Store)
	wsHandler := wshandlers.NewWebSocketHandlerCE(
		d,
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	)
	return &Server{
		wsServer: ws.NewServerCE(wsMiddle, wsHandler),
		httpServer: httpserver.NewServerCE(
			httpMiddle,
			apiKeyHandler,
			environmentHandler,
			groupHandler,
			hostInstanceHandler,
			organizationHandler,
			pageHandler,
			userHandler,
		),
	}
}

func (s *Server) router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(600) * time.Second))
	r.Use(cors.New(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			var pattern string

			switch config.Config.Env {
			case config.EnvProd:
				pattern = `^https://[a-zA-Z0-9-]+\.trysourcetool\.com$`
			case config.EnvStg:
				pattern = `^https://[a-zA-Z0-9-]+\.stg\.trysourcetool\.com$`
			case config.EnvLocal:
				pattern = `^(http://[a-zA-Z0-9-]+\.local\.trysourcetool\.com:\d+|http://localhost:\d+)$`
			default:
				return false
			}

			matched, err := regexp.MatchString(pattern, origin)
			if err != nil {
				return false
			}

			return matched
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           0,
		Debug:            !(config.Config.Env == config.EnvProd),
	}).Handler)

	if config.Config.Env == config.EnvLocal {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", "http://localhost:8080")),
		))
	}

	r.Mount("/ws", s.wsServer.Router())
	r.Mount("/api", s.httpServer.Router())

	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router().ServeHTTP(w, r)
}
