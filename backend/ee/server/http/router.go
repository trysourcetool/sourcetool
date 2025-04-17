package http

import (
	apikeySvc "github.com/trysourcetool/sourcetool/backend/ee/apikey/service"
	authSvc "github.com/trysourcetool/sourcetool/backend/ee/auth/service"
	environmentSvc "github.com/trysourcetool/sourcetool/backend/ee/environment/service"
	groupSvc "github.com/trysourcetool/sourcetool/backend/ee/group/service"
	hostinstanceSvc "github.com/trysourcetool/sourcetool/backend/ee/hostinstance/service"
	organizationSvc "github.com/trysourcetool/sourcetool/backend/ee/organization/service"
	pageSvc "github.com/trysourcetool/sourcetool/backend/ee/page/service"
	userSvc "github.com/trysourcetool/sourcetool/backend/ee/user/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/server/http"
	"github.com/trysourcetool/sourcetool/backend/server/http/handlers"
)

func NewRouter(d *infra.Dependency) *http.Router {
	middleware := NewMiddlewareEE(d.Store)
	apiKeyHandler := handlers.NewAPIKeyHandler(apikeySvc.NewAPIKeyServiceEE(d))
	authHandler := handlers.NewAuthHandler(authSvc.NewAuthServiceEE(d))
	environmentHandler := handlers.NewEnvironmentHandler(environmentSvc.NewEnvironmentServiceEE(d))
	groupHandler := handlers.NewGroupHandler(groupSvc.NewGroupServiceEE(d))
	hostInstanceHandler := handlers.NewHostInstanceHandler(hostinstanceSvc.NewHostInstanceServiceEE(d))
	organizationHandler := handlers.NewOrganizationHandler(organizationSvc.NewOrganizationServiceEE(d))
	pageHandler := handlers.NewPageHandler(pageSvc.NewPageServiceEE(d))
	userHandler := handlers.NewUserHandler(userSvc.NewUserServiceEE(d))
	return http.NewRouter(
		middleware,
		apiKeyHandler,
		authHandler,
		environmentHandler,
		groupHandler,
		hostInstanceHandler,
		organizationHandler,
		pageHandler,
		userHandler,
	)
}
