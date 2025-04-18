package user

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/user"
)

type repositoryEE struct {
	db      db.DB
	builder sq.StatementBuilderType
	*user.RepositoryCE
}

func NewRepositoryEE(db db.DB) *repositoryEE {
	return &repositoryEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		RepositoryCE: user.NewRepositoryCE(db),
	}
}
