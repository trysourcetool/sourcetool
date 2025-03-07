package http

import (
	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/server/http/handlers"
)

type Router struct {
	middleware   Middleware
	apikey       *handlers.APIKeyHandler
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

	r.Use(router.middleware.SetHTTPHeader)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signin", router.user.SignIn)
			r.Post("/signup/instructions", router.user.SendSignUpInstructions)
			r.Post("/signup", router.user.SignUp)
			r.Post("/oauth/google/authCodeUrl", router.user.GetGoogleAuthCodeURL)
			r.Post("/oauth/google/signin", router.user.SignInWithGoogle)
			r.Post("/oauth/google/signup", router.user.SignUpWithGoogle)
			r.Get("/oauth/google/callback", router.user.GoogleOAuthCallback)

			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthOrganization)
				r.Post("/saveAuth", router.user.SaveAuth)
				r.Post("/refreshToken", router.user.RefreshToken)
				r.Post("/invitations/signin", router.user.SignInInvitation)
				r.Post("/invitations/signup", router.user.SignUpInvitation)
				r.Post("/invitations/oauth/google/authCodeUrl", router.user.GetGoogleAuthCodeURLInvitation)
				r.Post("/invitations/oauth/google/signin", router.user.SignInWithGoogleInvitation)
				r.Post("/invitations/oauth/google/signup", router.user.SignUpWithGoogleInvitation)
			})

			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthUser)
				r.Get("/me", router.user.GetMe)
				r.Post("/obtainAuthToken", router.user.ObtainAuthToken)
			})

			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthUserWithOrganization)
				r.Get("/", router.user.List)
				r.Put("/", router.user.Update)
				r.Post("/sendUpdateEmailInstructions", router.user.SendUpdateEmailInstructions)
				r.Put("/email", router.user.UpdateEmail)
				r.Put("/password", router.user.UpdatePassword)
				r.Post("/invite", router.user.Invite)
				r.Post("/invitations/resend", router.user.ResendInvitation)
				r.Post("/signout", router.user.SignOut)
			})
		})

		r.Route("/organizations", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthUser)
				r.Post("/", router.organization.Create)
				r.Get("/checkSubdomainAvailability", router.organization.CheckSubdomainAvailability)
			})

			r.Group(func(r chi.Router) {
				r.Use(router.middleware.AuthUserWithOrganization)

				r.Route("/users", func(r chi.Router) {
					r.Put("/{userID}", router.organization.UpdateUser)
				})
			})
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
