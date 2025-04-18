package apikey

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
)

type serviceEE struct {
	*infra.Dependency
	*apikey.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: apikey.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer),
		),
	}
}
