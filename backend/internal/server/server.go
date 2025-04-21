package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	exceptionv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/exception/v1"
	websocketv1 "github.com/trysourcetool/sourcetool/backend/internal/pb/go/websocket/v1"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/pubsub"
	"github.com/trysourcetool/sourcetool/backend/internal/ws"
)

type Server struct {
	db        *postgres.DB
	pubsub    *pubsub.PubSub
	wsManager *ws.Manager
	checker   *permission.Checker
	upgrader  websocket.Upgrader
}

func New(
	db *postgres.DB,
	pubsub *pubsub.PubSub,
	wsManager *ws.Manager,
	checker *permission.Checker,
	upgrader websocket.Upgrader,
) *Server {
	return &Server{db, pubsub, wsManager, checker, upgrader}
}

func (s *Server) Install(router *chi.Mux) error {
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
				r.Get("/", s.errorHandler(s.listAPIKeys))
				r.Post("/", s.errorHandler(s.createAPIKey))

				r.Route("/{apiKeyID}", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.getAPIKey))
					r.Put("/", s.errorHandler(s.updateAPIKey))
					r.Delete("/", s.errorHandler(s.deleteAPIKey))
				})
			})

			r.Route("/auth", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Group(func(r chi.Router) {
						r.Use(s.authOrganizationIfSubdomainExists)

						r.Post("/magic/request", s.errorHandler(s.requestMagicLink))
						r.Post("/magic/authenticate", s.errorHandler(s.authenticateWithMagicLink))
						r.Post("/magic/register", s.errorHandler(s.registerWithMagicLink))
						r.Post("/invitations/magic/request", s.errorHandler(s.requestInvitationMagicLink))
						r.Post("/invitations/magic/authenticate", s.errorHandler(s.authenticateWithInvitationMagicLink))
						r.Post("/invitations/magic/register", s.errorHandler(s.registerWithInvitationMagicLink))

						r.Post("/google/request", s.errorHandler(s.requestGoogleAuthLink))
						r.Post("/google/authenticate", s.errorHandler(s.authenticateWithGoogle))
						r.Post("/google/register", s.errorHandler(s.registerWithGoogle))
						r.Post("/invitations/google/request", s.errorHandler(s.requestInvitationGoogleAuthLink))

						r.Post("/save", s.errorHandler(s.saveAuth))
						r.Post("/refresh", s.errorHandler(s.refreshToken))
					})

					r.Group(func(r chi.Router) {
						r.Use(s.authUserWithOrganizationIfSubdomainExists)
						r.Post("/token/obtain", s.errorHandler(s.obtainAuthToken))
					})

					r.Group(func(r chi.Router) {
						r.Use(s.authUserWithOrganization)
						r.Post("/logout", s.errorHandler(s.logout))
					})
				})

				r.Group(func(r chi.Router) {
					r.Use(s.authUserWithOrganizationIfSubdomainExists)
					r.Post("/token/obtain", s.errorHandler(s.obtainAuthToken))
				})

				r.Group(func(r chi.Router) {
					r.Use(s.authUserWithOrganization)
					r.Post("/logout", s.errorHandler(s.logout))
				})
			})

			r.Route("/environments", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.listEnvironments))
				r.Post("/", s.errorHandler(s.createEnvironment))

				r.Route("/{environmentID}", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.getEnvironment))
					r.Put("/", s.errorHandler(s.updateEnvironment))
					r.Delete("/", s.errorHandler(s.deleteEnvironment))
				})
			})

			r.Route("/groups", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.listGroups))
				r.Post("/", s.errorHandler(s.createGroup))

				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.getGroup))
					r.Put("/", s.errorHandler(s.updateGroup))
					r.Delete("/", s.errorHandler(s.deleteGroup))
				})
			})

			r.Route("/hostInstances", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/ping", s.errorHandler(s.pingHostInstance))
			})

			r.Route("/organizations", func(r chi.Router) {
				r.Use(s.authUser)
				r.Post("/", s.errorHandler(s.createOrganization))
				r.Get("/checkSubdomainAvailability", s.errorHandler(s.checkOrganizationSubdomainAvailability))
			})

			r.Route("/pages", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)
				r.Get("/", s.errorHandler(s.listPages))
			})

			r.Route("/users", func(r chi.Router) {
				r.Use(s.authUserWithOrganization)

				// Authenticated User
				r.Route("/me", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.getMe))
					r.Put("/", s.errorHandler(s.updateMe))
					r.Post("/email/instructions", s.errorHandler(s.sendUpdateMeEmailInstructions))
					r.Put("/email", s.errorHandler(s.updateMeEmail))
				})

				// Organization Users
				r.Get("/", s.errorHandler(s.listUsers))
				r.Route("/{userID}", func(r chi.Router) {
					r.Put("/", s.errorHandler(s.updateUser))
					r.Delete("/", s.errorHandler(s.deleteUser))
				})

				// Organization Invitations
				r.Route("/invitations", func(r chi.Router) {
					r.Post("/", s.errorHandler(s.createUserInvitations))
					r.Post("/{invitationID}/resend", s.errorHandler(s.resendUserInvitation))
				})
			})
		})
	})

	router.Route("/ws", func(r chi.Router) {
		r.Use(s.authWebSocketUser)
		r.Get("/", s.handleWebSocket)
	})

	staticDir := os.Getenv("STATIC_FILES_DIR")
	serveStaticFiles(router, staticDir)

	return nil
}

func (s *Server) errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			s.serveError(w, r, err)
		}
	}
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	currentUser := internal.CurrentUser(ctx)
	var email string
	if currentUser != nil {
		email = currentUser.Email
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
	currentUser := internal.CurrentUser(ctx)
	var email string
	if currentUser != nil {
		email = currentUser.Email
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
