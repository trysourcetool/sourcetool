package service

import (
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/user/service"
)

type userServiceEE struct {
	*infra.Dependency
	*service.UserServiceCE
}

func NewUserServiceEE(d *infra.Dependency) *userServiceEE {
	return &userServiceEE{
		Dependency: d,
		UserServiceCE: service.NewUserServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}
