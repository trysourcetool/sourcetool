package transport

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

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/app/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/app/auth"
	"github.com/trysourcetool/sourcetool/backend/internal/app/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/app/group"
	"github.com/trysourcetool/sourcetool/backend/internal/app/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/app/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/app/page"
	"github.com/trysourcetool/sourcetool/backend/internal/app/user"
	wsSvc "github.com/trysourcetool/sourcetool/backend/internal/app/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	v1 "github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1"
	v1handlers "github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/handlers"
	wstransport "github.com/trysourcetool/sourcetool/backend/internal/transport/ws"
	wshandlers "github.com/trysourcetool/sourcetool/backend/internal/transport/ws/handlers"
)

type Router struct {
	wsRouter   *wstransport.Router
	httpRouter *v1.Router
}

func NewRouter(d *infra.Dependency) *Router {
	httpMiddle := v1.NewMiddlewareCE(d.Repository)
	apiKeyHandler := v1handlers.NewAPIKeyHandler(apikey.NewServiceCE(d))
	authHandler := v1handlers.NewAuthHandler(auth.NewServiceCE(d))
	environmentHandler := v1handlers.NewEnvironmentHandler(environment.NewServiceCE(d))
	groupHandler := v1handlers.NewGroupHandler(group.NewServiceCE(d))
	hostInstanceHandler := v1handlers.NewHostInstanceHandler(hostinstance.NewServiceCE(d))
	organizationHandler := v1handlers.NewOrganizationHandler(organization.NewServiceCE(d))
	pageHandler := v1handlers.NewPageHandler(page.NewServiceCE(d))
	userHandler := v1handlers.NewUserHandler(user.NewServiceCE(d))

	wsMiddle := wstransport.NewMiddlewareCE(d.Repository)
	wsHandler := wshandlers.NewWebSocketHandler(
		websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		d.WSManager,
		wsSvc.NewServiceCE(d),
	)
	return &Router{
		wsRouter: wstransport.NewRouter(wsMiddle, wsHandler),
		httpRouter: v1.NewRouter(
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

func (r *Router) Build() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Duration(600) * time.Second))
	router.Use(cors.New(cors.Options{
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
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", "http://localhost:8080")),
		))
	}

	router.Mount("/ws", r.wsRouter.Build())
	router.Mount("/api", r.httpRouter.Build())

	staticDir := os.Getenv("STATIC_FILES_DIR")
	ServeStaticFiles(router, staticDir)

	return router
}
