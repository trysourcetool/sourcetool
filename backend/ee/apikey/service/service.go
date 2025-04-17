package service

import (
	"github.com/trysourcetool/sourcetool/backend/apikey/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type apikeyServiceEE struct {
	*infra.Dependency
	*service.APIKeyServiceCE
}

func NewAPIKeyServiceEE(d *infra.Dependency) *apikeyServiceEE {
	return &apikeyServiceEE{
		Dependency: d,
		APIKeyServiceCE: service.NewAPIKeyServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
