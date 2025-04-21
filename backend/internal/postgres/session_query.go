package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type SessionQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isSessionQuery()
}

type sessionByIDQuery struct{ id uuid.UUID }

func (q sessionByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`s."id"`: q.id})
}

func (sessionByIDQuery) isSessionQuery() {}

func SessionByID(id uuid.UUID) SessionQuery { return sessionByIDQuery{id: id} }
