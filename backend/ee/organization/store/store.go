package store

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization/store"
)

type organizationStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.OrganizationStoreCE
}

func NewOrganizationStoreEE(db infra.DB) *organizationStoreEE {
	return &organizationStoreEE{
		db:                  db,
		builder:             sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		OrganizationStoreCE: store.NewOrganizationStoreCE(db),
	}
}
