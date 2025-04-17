package apikey

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/apikey/store"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type apikeyStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.APIKeyStoreCE
}

func NewAPIKeyStoreEE(db infra.DB) *apikeyStoreEE {
	return &apikeyStoreEE{
		db:            db,
		builder:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		APIKeyStoreCE: store.NewAPIKeyStoreCE(db),
	}
}
