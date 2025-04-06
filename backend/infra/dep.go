package infra

type Dependency struct {
	Store  Store
	Mailer Mailer
	Memory Memory
}

func NewDependency(store Store, mailer Mailer) *Dependency {
	redisClient := NewRedisClientCE()
	memory := NewMemoryCE(redisClient)
	
	return &Dependency{
		Store:  store,
		Mailer: mailer,
		Memory: memory,
	}
}
