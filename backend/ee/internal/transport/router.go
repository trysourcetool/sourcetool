package transport

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/trysourcetool/sourcetool/backend/config"
	httpserver "github.com/trysourcetool/sourcetool/backend/ee/internal/transport/http/v1"
	"github.com/trysourcetool/sourcetool/backend/ee/internal/transport/ws"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	"github.com/trysourcetool/sourcetool/backend/internal/transport"
	cehttpserver "github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1"
	cews "github.com/trysourcetool/sourcetool/backend/internal/transport/ws"
)

type Router struct {
	wsRouter   *cews.Router
	httpRouter *cehttpserver.Router
}

func NewRouter(d *infra.Dependency) *Router {
	return &Router{
		wsRouter:   ws.NewRouter(d),
		httpRouter: httpserver.NewRouter(d),
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
			var pattern string

			switch config.Config.Env {
			case config.EnvProd:
				pattern = `^https://[a-zA-Z0-9-]+\.trysourcetool\.com$`
			case config.EnvStaging:
				pattern = `^https://[a-zA-Z0-9-]+\.staging\.trysourcetool\.com$`
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
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", "http://localhost:8080")),
		))
	}

	router.Mount("/ws", r.wsRouter.Build())
	router.Mount("/api", r.httpRouter.Build())

	staticDir := os.Getenv("STATIC_FILES_DIR")
	transport.ServeStaticFiles(router, staticDir)

	return router
}
