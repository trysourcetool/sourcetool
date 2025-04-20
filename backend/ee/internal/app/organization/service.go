package organization

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
)

type serviceEE struct {
	*port.Dependencies
	*organization.ServiceCE
}

func NewServiceEE(d *port.Dependencies) *serviceEE {
	return &serviceEE{
		Dependencies: d,
		ServiceCE: organization.NewServiceCE(
			port.NewDependencies(d.Repository, d.Mailer, d.PubSub, d.WSManager),
		),
	}
}
