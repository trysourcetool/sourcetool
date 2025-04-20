package environment

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
)

type serviceEE struct {
	*port.Dependencies
	*environment.ServiceCE
}

func NewServiceEE(d *port.Dependencies) *serviceEE {
	return &serviceEE{
		Dependencies: d,
		ServiceCE: environment.NewServiceCE(
			port.NewDependencies(d.Repository, d.Mailer, d.PubSub, d.WSManager),
		),
	}
}
