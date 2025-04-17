package page

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/page/store"
)

type pageStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.PageStoreCE
}

func NewPageStoreEE(db infra.DB) *pageStoreEE {
	return &pageStoreEE{
		db:          db,
		builder:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		PageStoreCE: store.NewPageStoreCE(db),
	}
}
