package organization

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
)

type serviceEE struct {
	*infra.Dependency
	*organization.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: organization.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer),
		),
	}
}
