package organization

import (
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization"
)

type serviceEE struct {
	*infra.Dependency
	*organization.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: organization.NewServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
