package user

import (
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type ServiceEE interface {
	user.ServiceCE
}

type serviceEE struct {
	*infra.Dependency
	*user.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: user.NewServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
