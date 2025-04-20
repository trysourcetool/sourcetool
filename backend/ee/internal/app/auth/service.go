package auth

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/auth"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
)

type serviceEE struct {
	*port.Dependencies
	*auth.ServiceCE
}

func NewServiceEE(d *port.Dependencies) *serviceEE {
	return &serviceEE{
		Dependencies: d,
		ServiceCE: auth.NewServiceCE(
			port.NewDependencies(d.Repository, d.Mailer, d.PubSub, d.WSManager),
		),
	}
}
