package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type SessionOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func SessionByID(id uuid.UUID) SessionOption {
	return sessionByIDOption{id: id}
}

type sessionByIDOption struct {
	id uuid.UUID
}

func (o sessionByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`s."id"`: o.id})
}
