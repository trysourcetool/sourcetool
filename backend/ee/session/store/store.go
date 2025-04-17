package store

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/session/store"
)

type sessionStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.SessionStoreCE
}

func NewSessionStoreEE(db infra.DB) *sessionStoreEE {
	return &sessionStoreEE{
		db:             db,
		builder:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		SessionStoreCE: store.NewSessionStoreCE(db),
	}
}
