package infra

import (
	"github.com/trysourcetool/sourcetool/backend/health"
)

type Dependency struct {
	Store  Store
	Mailer Mailer
	Health health.Service
}

func NewDependency(store Store, mailer Mailer) *Dependency {
	dep := &Dependency{
		Store:  store,
		Mailer: mailer,
	}
	
	dep.Health = health.NewServiceCE(dep)
	
	return dep
}
