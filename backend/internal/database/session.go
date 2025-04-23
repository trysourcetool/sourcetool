package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type SessionQuery interface{ isSessionQuery() }

type SessionByIDQuery struct{ ID uuid.UUID }

func (q SessionByIDQuery) isSessionQuery() {}

func SessionByID(id uuid.UUID) SessionQuery { return SessionByIDQuery{ID: id} }

type SessionStore interface {
	Get(ctx context.Context, queries ...SessionQuery) (*core.Session, error)
	Create(ctx context.Context, m *core.Session) error
	Delete(ctx context.Context, m *core.Session) error

	CreateHostInstance(ctx context.Context, m *core.SessionHostInstance) error
}
