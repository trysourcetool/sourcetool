package auth

import (
	"github.com/trysourcetool/sourcetool/backend/auth"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type serviceEE struct {
	*infra.Dependency
	*auth.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: auth.NewServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
