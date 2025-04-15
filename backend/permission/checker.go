package permission

import "github.com/trysourcetool/sourcetool/backend/infra"

type Checker struct {
	store infra.ModelStore
}

func NewChecker(store infra.ModelStore) *Checker {
	return &Checker{store}
}
