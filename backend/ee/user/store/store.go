package store

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/user/store"
)

type userStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.UserStoreCE
}

func NewUserStoreEE(db infra.DB) *userStoreEE {
	return &userStoreEE{
		db:          db,
		builder:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		UserStoreCE: store.NewUserStoreCE(db),
	}
}
