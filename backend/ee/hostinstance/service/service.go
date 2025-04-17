package service

import (
	"github.com/trysourcetool/sourcetool/backend/hostinstance/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type hostinstanceServiceEE struct {
	*infra.Dependency
	*service.HostInstanceServiceCE
}

func NewHostInstanceServiceEE(d *infra.Dependency) *hostinstanceServiceEE {
	return &hostinstanceServiceEE{
		Dependency: d,
		HostInstanceServiceCE: service.NewHostInstanceServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
