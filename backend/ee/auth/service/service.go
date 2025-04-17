package service

import (
	"github.com/trysourcetool/sourcetool/backend/auth/service"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type authServiceEE struct {
	*infra.Dependency
	*service.AuthServiceCE
}

func NewAuthServiceEE(d *infra.Dependency) *authServiceEE {
	return &authServiceEE{
		Dependency: d,
		AuthServiceCE: service.NewAuthServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
