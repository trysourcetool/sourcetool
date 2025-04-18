package v1

import (
	apikeySvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/apikey"
	authSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/auth"
	environmentSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/environment"
	groupSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/group"
	hostinstanceSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/hostinstance"
	organizationSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/organization"
	pageSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/page"
	userSvc "github.com/trysourcetool/sourcetool/backend/ee/internal/app/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	v1 "github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/handlers"
)

func NewRouter(d *infra.Dependency) *v1.Router {
	middleware := NewMiddlewareEE(d)
	apiKeyHandler := handlers.NewAPIKeyHandler(apikeySvc.NewServiceEE(d))
	authHandler := handlers.NewAuthHandler(authSvc.NewServiceEE(d))
	environmentHandler := handlers.NewEnvironmentHandler(environmentSvc.NewServiceEE(d))
	groupHandler := handlers.NewGroupHandler(groupSvc.NewServiceEE(d))
	hostInstanceHandler := handlers.NewHostInstanceHandler(hostinstanceSvc.NewServiceEE(d))
	organizationHandler := handlers.NewOrganizationHandler(organizationSvc.NewServiceEE(d))
	pageHandler := handlers.NewPageHandler(pageSvc.NewServiceEE(d))
	userHandler := handlers.NewUserHandler(userSvc.NewServiceEE(d))
	return v1.NewRouter(
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
