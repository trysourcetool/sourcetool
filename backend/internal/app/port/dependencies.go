package port

type Dependencies struct {
	Repository Repository
	Mailer     Mailer
	PubSub     PubSub
	WSManager  WSManager
}

func NewDependencies(repo Repository, mailer Mailer, pubsub PubSub, wsManager WSManager) *Dependencies {
	return &Dependencies{
		Repository: repo,
		Mailer:     mailer,
		PubSub:     pubsub,
		WSManager:  wsManager,
	}
}
