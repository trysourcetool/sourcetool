package server

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/trysourcetool/sourcetool/backend/ee/config"
	httpserver "github.com/trysourcetool/sourcetool/backend/ee/server/http"
	"github.com/trysourcetool/sourcetool/backend/ee/server/ws"
	"github.com/trysourcetool/sourcetool/backend/infra"
	cehttpserver "github.com/trysourcetool/sourcetool/backend/server/http"
	cews "github.com/trysourcetool/sourcetool/backend/server/ws"
)

type Server struct {
	wsRouter   *cews.Router
	httpRouter *cehttpserver.Router
}

func New(d *infra.Dependency) *Server {
	return &Server{
		wsRouter:   ws.NewRouter(d),
		httpRouter: httpserver.NewRouter(d),
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

	r.Mount("/ws", s.wsRouter.Build())
	r.Mount("/api", s.httpRouter.Build())

	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router().ServeHTTP(w, r)
}
