package environment

import (
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type serviceEE struct {
	*infra.Dependency
	*environment.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: environment.NewServiceCE(
			infra.NewDependency(d.Store, d.Signer, d.Mailer),
		),
	}
}
