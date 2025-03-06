package environment

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type storeEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*environment.StoreCE
}

func NewStoreEE(db infra.DB) *storeEE {
	return &storeEE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		StoreCE: environment.NewStoreCE(db),
	}
}
