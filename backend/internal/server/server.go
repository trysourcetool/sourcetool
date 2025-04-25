package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	exceptionv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/exception/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/pubsub"
	"github.com/trysourcetool/sourcetool/backend/internal/ws"
)

type Server struct {
	db        database.DB
	pubsub    *pubsub.PubSub
	wsManager *ws.Manager
	checker   *permission.Checker
	upgrader  websocket.Upgrader
}

func New(
	db database.DB,
	pubsub *pubsub.PubSub,
	wsManager *ws.Manager,
	checker *permission.Checker,
	upgrader websocket.Upgrader,
) *Server {
	return &Server{db, pubsub, wsManager, checker, upgrader}
}

func (s *Server) installDefaultMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Duration(600) * time.Second))
}

func (s *Server) installCORSMiddleware(router *chi.Mux) {
	router.Use(cors.New(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			if config.Config.IsCloudEdition {
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
			} else {
				// For self-hosted environments, we only need to check against the configured base URL
				normalizedOrigin := strings.TrimRight(origin, "/")
				normalizedBaseURL := strings.TrimRight(config.Config.BaseURL, "/")
				return normalizedOrigin == normalizedBaseURL
			}
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
}

func (s *Server) errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			s.serveError(w, r, err)
		}
	}
}

func (s *Server) installRESTHandlers(router *chi.Mux) {
	router.Route("/api", func(r chi.Router) {
		r.Use(s.setSubdomain)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		})

		r.Route("/v1", func(r chi.Router) {
			r.Route("/apiKeys", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.handleListAPIKeys))
				r.Post("/", s.errorHandler(s.handleCreateAPIKey))

				r.Route("/{apiKeyID}", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.handleGetAPIKey))
					r.Put("/", s.errorHandler(s.handleUpdateAPIKey))
					r.Delete("/", s.errorHandler(s.handleDeleteAPIKey))
				})
			})

			r.Route("/auth", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(s.authOrganizationIfSubdomainExists)

					r.Post("/magic/request", s.errorHandler(s.handleRequestMagicLink))
					r.Post("/magic/authenticate", s.errorHandler(s.handleAuthenticateWithMagicLink))
					r.Post("/magic/register", s.errorHandler(s.handleRegisterWithMagicLink))
					r.Post("/invitations/magic/request", s.errorHandler(s.handleRequestInvitationMagicLink))
					r.Post("/invitations/magic/authenticate", s.errorHandler(s.handleAuthenticateWithInvitationMagicLink))
					r.Post("/invitations/magic/register", s.errorHandler(s.handleRegisterWithInvitationMagicLink))

					r.Post("/google/request", s.errorHandler(s.handleRequestGoogleAuthLink))
					r.Post("/google/authenticate", s.errorHandler(s.handleAuthenticateWithGoogle))
					r.Post("/google/register", s.errorHandler(s.handleRegisterWithGoogle))
					r.Post("/invitations/google/request", s.errorHandler(s.handleRequestInvitationGoogleAuthLink))

					r.Post("/save", s.errorHandler(s.handleSaveAuth))
					r.Post("/refresh", s.errorHandler(s.handleRefreshToken))
				})

				r.Group(func(r chi.Router) {
					r.Use(s.authUserWithOrganizationIfSubdomainExists)
					r.Post("/token/obtain", s.errorHandler(s.handleObtainAuthToken))
				})

				r.Group(func(r chi.Router) {
					r.Use(s.authUserWithOrganization)
					r.Post("/logout", s.errorHandler(s.handleLogout))
				})
			})

			r.Route("/environments", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.handleListEnvironments))
				r.Post("/", s.errorHandler(s.handleCreateEnvironment))

				r.Route("/{environmentID}", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.handleGetEnvironment))
					r.Put("/", s.errorHandler(s.handleUpdateEnvironment))
					r.Delete("/", s.errorHandler(s.handleDeleteEnvironment))
				})
			})

			r.Route("/groups", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.handleListGroups))
				r.Post("/", s.errorHandler(s.handleCreateGroup))

				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.handleGetGroup))
					r.Put("/", s.errorHandler(s.handleUpdateGroup))
					r.Delete("/", s.errorHandler(s.handleDeleteGroup))
				})
			})

			r.Route("/hostInstances", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/ping", s.errorHandler(s.handlePingHostInstance))
			})

			r.Route("/organizations", func(r chi.Router) {
				r.Use(s.authUser)
				r.Post("/", s.errorHandler(s.handleCreateOrganization))
				r.Get("/checkSubdomainAvailability", s.errorHandler(s.handleCheckOrganizationSubdomainAvailability))
			})

			r.Route("/pages", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.handleListPages))
			})

			r.Route("/users", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)

				// Authenticated User
				r.Route("/me", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.handleGetMe))
					r.Put("/", s.errorHandler(s.handleUpdateMe))
					r.Post("/email/instructions", s.errorHandler(s.handleSendUpdateMeEmailInstructions))
					r.Put("/email", s.errorHandler(s.handleUpdateMeEmail))
				})

				// Organization Users
				r.Get("/", s.errorHandler(s.handleListUsers))
				r.Route("/{userID}", func(r chi.Router) {
					r.Put("/", s.errorHandler(s.handleUpdateUser))
					r.Delete("/", s.errorHandler(s.handleDeleteUser))
				})

				// Organization Invitations
				r.Route("/invitations", func(r chi.Router) {
					r.Post("/", s.errorHandler(s.handleCreateUserInvitations))
					r.Post("/{invitationID}/resend", s.errorHandler(s.handleResendUserInvitation))
				})
			})
		})
	})
}

func (s *Server) installWebSocketHandler(router *chi.Mux) {
	router.Route("/ws", func(r chi.Router) {
		r.Use(s.authWebSocketUser)
		r.Get("/", s.handleWebSocket)
	})
}

func (s *Server) installStaticHandler(router *chi.Mux) {
	staticDir := os.Getenv("STATIC_FILES_DIR")
	serveStaticFiles(router, staticDir)
}

func (s *Server) Install(router *chi.Mux) {
	s.installDefaultMiddlewares(router)
	s.installCORSMiddleware(router)
	s.installRESTHandlers(router)
	s.installWebSocketHandler(router)
	s.installStaticHandler(router)
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	var email string
	if ctxUser != nil {
		email = ctxUser.Email
	}

	v, ok := err.(*errdefs.Error)
	if !ok {
		logger.Logger.Error(
			err.Error(),
			zap.Stack("stack_trace"),
			zap.String("email", email),
			zap.String("cause", "application"),
		)

		s.renderJSON(
			w,
			http.StatusInternalServerError,
			errdefs.ErrInternal(err),
		)
		return
	}

	fields := []zap.Field{
		zap.String("email", email),
		zap.String("error_stacktrace", strings.Join(v.StackTrace(), "\n")),
	}

	switch {
	case v.Status >= 500:
		fields = append(fields, zap.String("cause", "application"))
		logger.Logger.Error(err.Error(), fields...)
	case v.Status >= 402, v.Status == 400:
		fields = append(fields, zap.String("cause", "user"))
		logger.Logger.Error(err.Error(), fields...)
	default:
		fields = append(fields, zap.String("cause", "internal_info"))
		logger.Logger.Warn(err.Error(), fields...)
	}

	s.renderJSON(w, v.Status, v)
}

func (s *Server) renderJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendWebSocketMessage(conn *websocket.Conn, msg *websocketv1.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendErrWebSocketMessage(ctx context.Context, conn *websocket.Conn, id string, err error) {
	ctxUser := internal.ContextUser(ctx)
	var email string
	if ctxUser != nil {
		email = ctxUser.Email
	}

	e, ok := err.(*errdefs.Error)
	if !ok {
		logger.Logger.Error(
			err.Error(),
			zap.Stack("stack_trace"),
			zap.String("email", email),
			zap.String("cause", "application"),
		)

		v := errdefs.ErrInternal(err)
		e, _ = v.(*errdefs.Error)
	} else {
		fields := []zap.Field{
			zap.String("email", email),
			zap.String("error_stacktrace", strings.Join(e.StackTrace(), "\n")),
		}

		switch {
		case e.Status >= 500:
			fields = append(fields, zap.String("cause", "application"))
			logger.Logger.Error(err.Error(), fields...)
		case e.Status >= 402, e.Status == 400:
			fields = append(fields, zap.String("cause", "user"))
			logger.Logger.Error(err.Error(), fields...)
		default:
			fields = append(fields, zap.String("cause", "internal_info"))
			logger.Logger.Warn(err.Error(), fields...)
		}
	}

	msg := &websocketv1.Message{
		Id: id,
		Type: &websocketv1.Message_Exception{
			Exception: &exceptionv1.Exception{
				Title:      e.Title,
				Message:    e.Detail,
				StackTrace: e.StackTrace(),
			},
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Failed to marshal WS error message", zap.Error(err))
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.Logger.Error("Failed to write WS error message", zap.Error(err))
		return
	}
}

type statusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
