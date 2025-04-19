package infra

import (
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/email"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/pubsub"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/ws"
)

type Dependency struct {
	Repository db.Repository
	Mailer     email.Mailer
	PubSub     pubsub.PubSub
	WSManager  ws.Manager
}

func NewDependency(repo db.Repository, mailer email.Mailer, pubsub pubsub.PubSub, wsManager ws.Manager) *Dependency {
	return &Dependency{
		Repository: repo,
		Mailer:     mailer,
		PubSub:     pubsub,
		WSManager:  wsManager,
	}
}
