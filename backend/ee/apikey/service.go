package apikey

import (
	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type serviceEE struct {
	*infra.Dependency
	*apikey.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: apikey.NewServiceCE(
			infra.NewDependency(d.Store, d.Signer, d.Mailer),
		),
	}
}
