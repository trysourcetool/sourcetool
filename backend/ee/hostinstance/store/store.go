package store

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/hostinstance/store"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type hostinstanceStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.HostInstanceStoreCE
}

func NewHostInstanceStoreEE(db infra.DB) *hostinstanceStoreEE {
	return &hostinstanceStoreEE{
		db:                  db,
		builder:             sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		HostInstanceStoreCE: store.NewHostInstanceStoreCE(db),
	}
}
