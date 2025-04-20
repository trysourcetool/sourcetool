package session

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/session"
)

type repositoryEE struct {
	db      db.DB
	builder sq.StatementBuilderType
	*session.RepositoryCE
}

func NewRepositoryEE(db db.DB) *repositoryEE {
	return &repositoryEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		RepositoryCE: session.NewRepositoryCE(db),
	}
}
