package apikey

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/db"
)

type repositoryEE struct {
	db      db.DB
	builder sq.StatementBuilderType
	*apikey.RepositoryCE
}

func NewRepositoryEE(db db.DB) *repositoryEE {
	return &repositoryEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		RepositoryCE: apikey.NewRepositoryCE(db),
	}
}
