package hostinstance

import (
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type serviceEE struct {
	*infra.Dependency
	*hostinstance.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: hostinstance.NewServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
