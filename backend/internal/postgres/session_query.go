package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type SessionQuery interface{ isSessionQuery() }

type SessionByIDQuery struct{ ID uuid.UUID }

func (q SessionByIDQuery) isSessionQuery() {}

func SessionByID(id uuid.UUID) SessionQuery { return SessionByIDQuery{ID: id} }

func applySessionQueries(b sq.SelectBuilder, queries ...SessionQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case SessionByIDQuery:
			b = b.Where(sq.Eq{`s."id"`: q.ID})
		}
	}

	return b
}
