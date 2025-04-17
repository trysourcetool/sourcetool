package service

import (
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization/service"
)

type organizationServiceEE struct {
	*infra.Dependency
	*service.OrganizationServiceCE
}

func NewOrganizationServiceEE(d *infra.Dependency) *organizationServiceEE {
	return &organizationServiceEE{
		Dependency: d,
		OrganizationServiceCE: service.NewOrganizationServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
