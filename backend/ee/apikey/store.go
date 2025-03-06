package apikey

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type storeEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*apikey.StoreCE
}

func NewStoreEE(db infra.DB) *storeEE {
	return &storeEE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		StoreCE: apikey.NewStoreCE(db),
	}
}
