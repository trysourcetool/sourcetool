package session

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type RepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isRepositoryOption()
}

func ByID(id uuid.UUID) RepositoryOption {
	return byIDOption{id: id}
}

type byIDOption struct {
	id uuid.UUID
}

func (o byIDOption) isRepositoryOption() {}

func (o byIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`s."id"`: o.id})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*Session, error)
	Create(context.Context, *Session) error
	Delete(context.Context, *Session) error
}
