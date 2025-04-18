package permission

import "github.com/trysourcetool/sourcetool/backend/internal/infra/db"

type Checker struct {
	repo db.Repository
}

func NewChecker(repo db.Repository) *Checker {
	return &Checker{repo}
}
