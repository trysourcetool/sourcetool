package hostinstance

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type storeEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*hostinstance.StoreCE
}

func NewStoreEE(db infra.DB) *storeEE {
	return &storeEE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		StoreCE: hostinstance.NewStoreCE(db),
	}
}
