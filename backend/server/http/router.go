package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/server/http/handlers"
)

type Router struct {
	middleware   Middleware
	apikey       *handlers.APIKeyHandler
	auth         *handlers.AuthHandler
	environment  *handlers.EnvironmentHandler
	group        *handlers.GroupHandler
	hostInstance *handlers.HostInstanceHandler
	organization *handlers.OrganizationHandler
	page         *handlers.PageHandler
	user         *handlers.UserHandler
}

func NewRouter(
	middleware Middleware,
	apiKeyHandler *handlers.APIKeyHandler,
	authHandler *handlers.AuthHandler,
	environmentHandler *handlers.EnvironmentHandler,
	groupHandler *handlers.GroupHandler,
	hostInstanceHandler *handlers.HostInstanceHandler,
	organizationHandler *handlers.OrganizationHandler,
	pageHandler *handlers.PageHandler,
	userHandler *handlers.UserHandler,
) *Router {
	return &Router{
		middleware:   middleware,
		apikey:       apiKeyHandler,
		auth:         authHandler,
		environment:  environmentHandler,
		group:        groupHandler,
		hostInstance: hostInstanceHandler,
		organization: organizationHandler,
		page:         pageHandler,
		user:         userHandler,
	}
}

func (router *Router) Build() chi.Router {
	r := chi.NewRouter()
	r.Use(router.middleware.SetSubdomain)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthOrganizationIfSubdomainExists)

				r.Post("/magic/request", router.auth.RequestMagicLink)
				r.Post("/magic/authenticate", router.auth.AuthenticateWithMagicLink)
				r.Post("/magic/register", router.auth.RegisterWithMagicLink)
				r.Post("/invitations/magic/request", router.auth.RequestInvitationMagicLink)
				r.Post("/invitations/magic/authenticate", router.auth.AuthenticateWithInvitationMagicLink)
				r.Post("/invitations/magic/register", router.auth.RegisterWithInvitationMagicLink)

				r.Post("/google/request", router.auth.RequestGoogleAuthLink)
				r.Post("/google/authenticate", router.auth.AuthenticateWithGoogle)
				r.Post("/google/register", router.auth.RegisterWithGoogle)
				r.Post("/invitations/google/request", router.auth.RequestInvitationGoogleAuthLink)

				r.Post("/save", router.auth.Save)
				r.Post("/refresh", router.auth.RefreshToken)
			})

			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthUserWithOrganizationIfSubdomainExists)
				r.Post("/token/obtain", router.auth.ObtainAuthToken)
			})

			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthUserWithOrganization)
				r.Post("/logout", router.auth.Logout)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(router.middleware.AuthUserWithOrganization)

			// Authenticated User
			r.Get("/me", router.user.GetMe)
			r.Put("/me", router.user.UpdateMe)
			r.Post("/me/email/instructions", router.user.SendUpdateMeEmailInstructions)
			r.Put("/me/email", router.user.UpdateMeEmail)

			// Organization Users
			r.Get("/", router.user.List)
			r.Put("/{userID}", router.user.Update)
			r.Delete("/{userID}", router.user.Delete)

			// Organization Invitations
			r.Post("/invitations", router.user.CreateUserInvitations)
			r.Post("/invitations/{invitationID}/resend", router.user.ResendUserInvitation)
		})

		r.Route("/organizations", func(r chi.Router) {
			r.Use(router.middleware.AuthUser)
			r.Post("/", router.organization.Create)
			r.Get("/checkSubdomainAvailability", router.organization.CheckSubdomainAvailability)
		})

		r.Route("/environments", func(r chi.Router) {
			r.Use(router.middleware.AuthUserWithOrganization)
			r.Get("/{environmentID}", router.environment.Get)
			r.Get("/", router.environment.List)
			r.Post("/", router.environment.Create)
			r.Put("/{environmentID}", router.environment.Update)
			r.Delete("/{environmentID}", router.environment.Delete)
		})

		r.Route("/groups", func(r chi.Router) {
			r.Use(router.middleware.AuthUserWithOrganization)
			r.Get("/{groupID}", router.group.Get)
			r.Get("/", router.group.List)
			r.Post("/", router.group.Create)
			r.Put("/{groupID}", router.group.Update)
			r.Delete("/{groupID}", router.group.Delete)
		})

		r.Route("/apiKeys", func(r chi.Router) {
			r.Use(router.middleware.AuthUserWithOrganization)
			r.Get("/{apiKeyID}", router.apikey.Get)
			r.Get("/", router.apikey.List)
			r.Post("/", router.apikey.Create)
			r.Put("/{apiKeyID}", router.apikey.Update)
			r.Delete("/{apiKeyID}", router.apikey.Delete)
		})

		r.Route("/pages", func(r chi.Router) {
			r.Use(router.middleware.AuthUserWithOrganization)
			r.Get("/", router.page.List)
		})

		r.Route("/hostInstances", func(r chi.Router) {
			r.Use(router.middleware.AuthUserWithOrganization)
			r.Get("/ping", router.hostInstance.Ping)
		})
	})

	return r
}
