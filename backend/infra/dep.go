package infra

type Dependency struct {
	Store  Store
	Mailer Mailer
}

func NewDependency(store Store, mailer Mailer) *Dependency {
	return &Dependency{
		Store:  store,
		Mailer: mailer,
	}
}
