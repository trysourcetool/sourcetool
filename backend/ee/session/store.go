package session

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/session"
)

type storeEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*session.StoreCE
}

func NewStoreEE(db infra.DB) *storeEE {
	return &storeEE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		StoreCE: session.NewStoreCE(db),
	}
}
