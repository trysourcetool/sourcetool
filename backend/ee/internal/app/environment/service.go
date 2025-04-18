package environment

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
)

type serviceEE struct {
	*infra.Dependency
	*environment.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: environment.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer, d.PubSub),
		),
	}
}
