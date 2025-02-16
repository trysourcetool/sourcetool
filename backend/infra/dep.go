package infra

type Dependency struct {
	Store  Store
	Signer Signer
	Mailer Mailer
}

func NewDependency(store Store, signer Signer, mailer Mailer) *Dependency {
	return &Dependency{
		Store:  store,
		Signer: signer,
		Mailer: mailer,
	}
}
