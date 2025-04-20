package page

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/page"
)

type repositoryEE struct {
	db      db.DB
	builder sq.StatementBuilderType
	*page.RepositoryCE
}

func NewRepositoryEE(db db.DB) *repositoryEE {
	return &repositoryEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		RepositoryCE: page.NewRepositoryCE(db),
	}
}
