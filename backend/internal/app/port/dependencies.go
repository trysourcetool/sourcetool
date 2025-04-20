package port

type Dependencies struct {
	Repository Repository
	Mailer     Mailer
	PubSub     PubSub
	WSManager  WSManager
}

// NewDependency returns the old Dependency struct (deprecated, use NewDependencies).
func NewDependencies(repo Repository, mailer Mailer, pubsub PubSub, wsManager WSManager) *Dependencies {
	return &Dependencies{
		Repository: repo,
		Mailer:     mailer,
		PubSub:     pubsub,
		WSManager:  wsManager,
	}
}
