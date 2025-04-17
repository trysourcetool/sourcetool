package environment

import (
	"github.com/trysourcetool/sourcetool/backend/environment/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type environmentServiceEE struct {
	*infra.Dependency
	*service.EnvironmentServiceCE
}

func NewEnvironmentServiceEE(d *infra.Dependency) *environmentServiceEE {
	return &environmentServiceEE{
		Dependency: d,
		EnvironmentServiceCE: service.NewEnvironmentServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
