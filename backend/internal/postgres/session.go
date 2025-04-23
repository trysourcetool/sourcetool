package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.SessionStore = (*sessionStore)(nil)

type sessionStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newSessionStore(db internal.DB) *sessionStore {
	return &sessionStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *sessionStore) Get(ctx context.Context, queries ...database.SessionQuery) (*core.Session, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Session{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrSessionNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *sessionStore) applyQueries(b sq.SelectBuilder, queries ...database.SessionQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.SessionByIDQuery:
			b = b.Where(sq.Eq{`s."id"`: q.ID})
		}
	}

	return b
}

func (s *sessionStore) buildQuery(ctx context.Context, queries ...database.SessionQuery) (string, []any, error) {
	q := s.builder.Select(
		`s."id"`,
		`s."organization_id"`,
		`s."user_id"`,
		`s."environment_id"`,
		`s."created_at"`,
		`s."updated_at"`,
	).
		From(`"session" s`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *sessionStore) Create(ctx context.Context, m *core.Session) error {
	if _, err := s.builder.
		Insert(`"session"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"user_id"`,
			`"environment_id"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.UserID,
			m.EnvironmentID,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *sessionStore) Delete(ctx context.Context, m *core.Session) error {
	if _, err := s.builder.
		Delete(`"session"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *sessionStore) CreateHostInstance(ctx context.Context, m *core.SessionHostInstance) error {
	if _, err := s.builder.
		Insert(`"session_host_instance"`).
		Columns(
			`"id"`,
			`"session_id"`,
			`"host_instance_id"`,
		).
		Values(
			m.ID,
			m.SessionID,
			m.HostInstanceID,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
