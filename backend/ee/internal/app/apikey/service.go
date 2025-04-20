package apikey

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
)

type serviceEE struct {
	*port.Dependencies
	*apikey.ServiceCE
}

func NewServiceEE(d *port.Dependencies) *serviceEE {
	return &serviceEE{
		Dependencies: d,
		ServiceCE: apikey.NewServiceCE(
			port.NewDependencies(d.Repository, d.Mailer, d.PubSub, d.WSManager),
		),
	}
}
