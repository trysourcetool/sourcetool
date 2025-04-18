package infra

import (
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/email"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/pubsub"
)

type Dependency struct {
	Repository db.Repository
	Mailer     email.Mailer
	PubSub     pubsub.PubSub
}

func NewDependency(repo db.Repository, mailer email.Mailer, pubsub pubsub.PubSub) *Dependency {
	return &Dependency{
		Repository: repo,
		Mailer:     mailer,
		PubSub:     pubsub,
	}
}
