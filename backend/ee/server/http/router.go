package http

import (
	"github.com/trysourcetool/sourcetool/backend/ee/apikey"
	"github.com/trysourcetool/sourcetool/backend/ee/environment"
	"github.com/trysourcetool/sourcetool/backend/ee/group"
	"github.com/trysourcetool/sourcetool/backend/ee/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/ee/organization"
	"github.com/trysourcetool/sourcetool/backend/ee/page"
	"github.com/trysourcetool/sourcetool/backend/ee/user"
	"github.com/trysourcetool/sourcetool/backend/health"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/server/http"
	"github.com/trysourcetool/sourcetool/backend/server/http/handlers"
)

func NewRouter(d *infra.Dependency) *http.Router {
	middleware := NewMiddlewareEE(d.Store)
	apiKeyHandler := handlers.NewAPIKeyHandler(apikey.NewServiceEE(d))
	environmentHandler := handlers.NewEnvironmentHandler(environment.NewServiceEE(d))
	groupHandler := handlers.NewGroupHandler(group.NewServiceEE(d))
	hostInstanceHandler := handlers.NewHostInstanceHandler(hostinstance.NewServiceEE(d))
	organizationHandler := handlers.NewOrganizationHandler(organization.NewServiceEE(d))
	pageHandler := handlers.NewPageHandler(page.NewServiceEE(d))
	userHandler := handlers.NewUserHandler(user.NewServiceEE(d))
	healthService := health.NewServiceCE(d)
	healthHandler := handlers.NewHealthHandler(healthService)
	
	return http.NewRouter(
		middleware,
		apiKeyHandler,
		environmentHandler,
		groupHandler,
		hostInstanceHandler,
		organizationHandler,
		pageHandler,
		userHandler,
		healthHandler,
	)
}
