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

	apikeySvc "github.com/trysourcetool/sourcetool/backend/apikey/service"
	authSvc "github.com/trysourcetool/sourcetool/backend/auth/service"
	"github.com/trysourcetool/sourcetool/backend/config"
	environmentSvc "github.com/trysourcetool/sourcetool/backend/environment/service"
	groupSvc "github.com/trysourcetool/sourcetool/backend/group/service"
	hostinstanceSvc "github.com/trysourcetool/sourcetool/backend/hostinstance/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
	organizationSvc "github.com/trysourcetool/sourcetool/backend/organization/service"
	pageSvc "github.com/trysourcetool/sourcetool/backend/page/service"
	httpserver "github.com/trysourcetool/sourcetool/backend/server/http"
	httphandlers "github.com/trysourcetool/sourcetool/backend/server/http/handlers"
	wsserver "github.com/trysourcetool/sourcetool/backend/server/ws"
	wshandlers "github.com/trysourcetool/sourcetool/backend/server/ws/handlers"
	userSvc "github.com/trysourcetool/sourcetool/backend/user/service"
	wsSvc "github.com/trysourcetool/sourcetool/backend/ws/service"
)

type Server struct {
	wsRouter   *wsserver.Router
	httpRouter *httpserver.Router
}

func New(d *infra.Dependency) *Server {
	httpMiddle := httpserver.NewMiddlewareCE(d.Store)
	apiKeyHandler := httphandlers.NewAPIKeyHandler(apikeySvc.NewAPIKeyServiceCE(d))
	authHandler := httphandlers.NewAuthHandler(authSvc.NewAuthServiceCE(d))
	environmentHandler := httphandlers.NewEnvironmentHandler(environmentSvc.NewEnvironmentServiceCE(d))
	groupHandler := httphandlers.NewGroupHandler(groupSvc.NewGroupServiceCE(d))
	hostInstanceHandler := httphandlers.NewHostInstanceHandler(hostinstanceSvc.NewHostInstanceServiceCE(d))
	organizationHandler := httphandlers.NewOrganizationHandler(organizationSvc.NewOrganizationServiceCE(d))
	pageHandler := httphandlers.NewPageHandler(pageSvc.NewPageServiceCE(d))
	userHandler := httphandlers.NewUserHandler(userSvc.NewUserServiceCE(d))

	wsMiddle := wsserver.NewMiddlewareCE(d.Store)
	wsHandler := wshandlers.NewWebSocketHandler(
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		wsSvc.NewWebSocketServiceCE(d),
	)
	return &Server{
		wsRouter: wsserver.NewRouter(wsMiddle, wsHandler),
		httpRouter: httpserver.NewRouter(
			httpMiddle,
			apiKeyHandler,
			authHandler,
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

	staticDir := os.Getenv("STATIC_FILES_DIR")
	ServeStaticFiles(r, staticDir)

	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router().ServeHTTP(w, r)
}
