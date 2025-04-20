package session

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Query interface{ isQuery() }

type ByIDQuery struct{ ID uuid.UUID }

func (q ByIDQuery) isQuery() {}

func ByID(id uuid.UUID) Query { return ByIDQuery{ID: id} }

type Repository interface {
	Get(context.Context, ...Query) (*Session, error)
	Create(context.Context, *Session) error
	Delete(context.Context, *Session) error
}
