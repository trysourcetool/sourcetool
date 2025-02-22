package http

import (
	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/server/http/handlers"
)

type ServerCE struct {
	middleware   MiddlewareCE
	apikey       handlers.APIKeyHandlerCE
	environment  handlers.EnvironmentHandlerCE
	group        handlers.GroupHandlerCE
	hostInstance handlers.HostInstanceHandlerCE
	organization handlers.OrganizationHandlerCE
	page         handlers.PageHandlerCE
	user         handlers.UserHandlerCE
}

func NewServerCE(
	middleware MiddlewareCE,
	apiKeyHandler handlers.APIKeyHandlerCE,
	environmentHandler handlers.EnvironmentHandlerCE,
	groupHandler handlers.GroupHandlerCE,
	hostInstanceHandler handlers.HostInstanceHandlerCE,
	organizationHandler handlers.OrganizationHandlerCE,
	pageHandler handlers.PageHandlerCE,
	userHandler handlers.UserHandlerCE,
) *ServerCE {
	return &ServerCE{
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

func (s *ServerCE) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(s.middleware.SetHTTPHeader)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signin", s.user.SignIn)
			r.Post("/signup/instructions", s.user.SendSignUpInstructions)
			r.Post("/signup", s.user.SignUp)
			r.Post("/oauth/google/authCodeUrl", s.user.GetGoogleAuthCodeURL)
			r.Post("/oauth/google/signin", s.user.SignInWithGoogle)
			r.Post("/oauth/google/signup", s.user.SignUpWithGoogle)
			r.Get("/oauth/google/callback", s.user.GoogleOAuthCallback)

			r.Group(func(r chi.Router) {
				r.Use(s.middleware.AuthOrganization)
				r.Post("/saveAuth", s.user.SaveAuth)
				r.Post("/refreshToken", s.user.RefreshToken)
				r.Post("/invitations/signin", s.user.SignInInvitation)
				r.Post("/invitations/signup", s.user.SignUpInvitation)
				r.Post("/invitations/oauth/google/authCodeUrl", s.user.GetGoogleAuthCodeURLInvitation)
				r.Post("/invitations/oauth/google/signin", s.user.SignInWithGoogleInvitation)
				r.Post("/invitations/oauth/google/signup", s.user.SignUpWithGoogleInvitation)
			})

			r.Group(func(r chi.Router) {
				r.Use(s.middleware.AuthUser)
				r.Get("/me", s.user.GetMe)
				r.Post("/obtainAuthToken", s.user.ObtainAuthToken)
			})

			r.Group(func(r chi.Router) {
				r.Use(s.middleware.AuthUserWithOrganization)
				r.Get("/", s.user.List)
				r.Put("/", s.user.Update)
				r.Post("/sendUpdateEmailInstructions", s.user.SendUpdateEmailInstructions)
				r.Put("/email", s.user.UpdateEmail)
				r.Put("/password", s.user.UpdatePassword)
				r.Post("/invite", s.user.Invite)
				r.Post("/signout", s.user.SignOut)
			})
		})

		r.Route("/organizations", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(s.middleware.AuthUser)
				r.Post("/", s.organization.Create)
				r.Get("/checkSubdomainAvailability", s.organization.CheckSubdomainAvailability)
			})

			r.Group(func(r chi.Router) {
				r.Use(s.middleware.AuthUserWithOrganization)

				r.Route("/users", func(r chi.Router) {
					r.Put("/{userID}", s.organization.UpdateUser)
				})
			})
		})

		r.Route("/environments", func(r chi.Router) {
			r.Use(s.middleware.AuthUserWithOrganization)
			r.Get("/{environmentID}", s.environment.Get)
			r.Get("/", s.environment.List)
			r.Post("/", s.environment.Create)
			r.Put("/{environmentID}", s.environment.Update)
			r.Delete("/{environmentID}", s.environment.Delete)
		})

		r.Route("/groups", func(r chi.Router) {
			r.Use(s.middleware.AuthUserWithOrganization)
			r.Get("/{groupID}", s.group.Get)
			r.Get("/", s.group.List)
			r.Post("/", s.group.Create)
			r.Put("/{groupID}", s.group.Update)
			r.Delete("/{groupID}", s.group.Delete)
		})

		r.Route("/apiKeys", func(r chi.Router) {
			r.Use(s.middleware.AuthUserWithOrganization)
			r.Get("/{apiKeyID}", s.apikey.Get)
			r.Get("/", s.apikey.List)
			r.Post("/", s.apikey.Create)
			r.Put("/{apiKeyID}", s.apikey.Update)
			r.Delete("/{apiKeyID}", s.apikey.Delete)
		})

		r.Route("/pages", func(r chi.Router) {
			r.Use(s.middleware.AuthUserWithOrganization)
			r.Get("/", s.page.List)
		})

		r.Route("/hostInstances", func(r chi.Router) {
			r.Use(s.middleware.AuthUserWithOrganization)
			r.Get("/ping", s.hostInstance.Ping)
		})
	})

	return r
}
