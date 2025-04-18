package auth

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/auth"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
)

type serviceEE struct {
	*infra.Dependency
	*auth.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: auth.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer),
		),
	}
}
