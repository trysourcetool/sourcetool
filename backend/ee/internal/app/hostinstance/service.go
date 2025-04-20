package hostinstance

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
)

type serviceEE struct {
	*port.Dependencies
	*hostinstance.ServiceCE
}

func NewServiceEE(d *port.Dependencies) *serviceEE {
	return &serviceEE{
		Dependencies: d,
		ServiceCE: hostinstance.NewServiceCE(
			port.NewDependencies(d.Repository, d.Mailer, d.PubSub, d.WSManager),
		),
	}
}
