package organization

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/db"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/organization"
)

type repositoryEE struct {
	db      db.DB
	builder sq.StatementBuilderType
	*organization.RepositoryCE
}

func NewRepositoryEE(db db.DB) *repositoryEE {
	return &repositoryEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		RepositoryCE: organization.NewRepositoryCE(db),
	}
}
