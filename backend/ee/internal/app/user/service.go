package user

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
)

type serviceEE struct {
	*infra.Dependency
	*user.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: user.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer, d.PubSub, d.WSManager),
		),
	}
}
