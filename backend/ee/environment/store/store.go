package environment

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/environment/store"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type environmentStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.EnvironmentStoreCE
}

func NewEnvironmentStoreEE(db infra.DB) *environmentStoreEE {
	return &environmentStoreEE{
		db:                 db,
		builder:            sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		EnvironmentStoreCE: store.NewEnvironmentStoreCE(db),
	}
}
