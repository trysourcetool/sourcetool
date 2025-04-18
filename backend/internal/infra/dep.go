package infra

import (
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/email"
)

type Dependency struct {
	Repository db.Repository
	Mailer     email.Mailer
}

func NewDependency(repo db.Repository, mailer email.Mailer) *Dependency {
	return &Dependency{
		Repository: repo,
		Mailer:     mailer,
	}
}
