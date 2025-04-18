package hostinstance

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
)

type serviceEE struct {
	*infra.Dependency
	*hostinstance.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: hostinstance.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer),
		),
	}
}
