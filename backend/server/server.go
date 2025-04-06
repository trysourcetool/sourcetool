package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
	"github.com/trysourcetool/sourcetool/backend/health"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/postgres"
	httpserver "github.com/trysourcetool/sourcetool/backend/server/http"
	httphandlers "github.com/trysourcetool/sourcetool/backend/server/http/handlers"
	wsserver "github.com/trysourcetool/sourcetool/backend/server/ws"
	wshandlers "github.com/trysourcetool/sourcetool/backend/server/ws/handlers"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/ws"
)

type Server struct {
	wsRouter   *wsserver.Router
	httpRouter *httpserver.Router
}

func New(d *infra.Dependency) *Server {
	httpMiddle := httpserver.NewMiddlewareCE(d.Store)
	apiKeyHandler := httphandlers.NewAPIKeyHandler(apikey.NewServiceCE(d))
	environmentHandler := httphandlers.NewEnvironmentHandler(environment.NewServiceCE(d))
	groupHandler := httphandlers.NewGroupHandler(group.NewServiceCE(d))
	hostInstanceHandler := httphandlers.NewHostInstanceHandler(hostinstance.NewServiceCE(d))
	organizationHandler := httphandlers.NewOrganizationHandler(organization.NewServiceCE(d))
	pageHandler := httphandlers.NewPageHandler(page.NewServiceCE(d))
	userHandler := httphandlers.NewUserHandler(user.NewServiceCE(d))

	wsMiddle := wsserver.NewMiddlewareCE(d.Store)
	wsHandler := wshandlers.NewWebSocketHandler(
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		ws.NewServiceCE(d),
	)
	return &Server{
		wsRouter: wsserver.NewRouter(wsMiddle, wsHandler),
		httpRouter: httpserver.NewRouter(
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
			// For self-hosted environments, we only need to check against the configured base URL
			normalizedOrigin := strings.TrimRight(origin, "/")
			normalizedBaseURL := strings.TrimRight(config.Config.BaseURL, "/")
			return normalizedOrigin == normalizedBaseURL
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

	r.Mount("/ws", s.wsRouter.Build())
	r.Mount("/api", s.httpRouter.Build())

	db, err := postgres.New()
	if err == nil {
		healthService := health.NewServiceCE(db)
		r.Get("/health", httphandlers.NewHealthHandler(healthService).Check)
	}

	staticDir := os.Getenv("STATIC_FILES_DIR")
	ServeStaticFiles(r, staticDir)

	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router().ServeHTTP(w, r)
}
