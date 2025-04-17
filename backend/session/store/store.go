package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/session"
)

type SessionStoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewSessionStoreCE(db infra.DB) *SessionStoreCE {
	return &SessionStoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *SessionStoreCE) Get(ctx context.Context, opts ...session.StoreOption) (*session.Session, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := session.Session{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrSessionNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *SessionStoreCE) buildQuery(ctx context.Context, opts ...session.StoreOption) (string, []any, error) {
	q := s.builder.Select(
		`s."id"`,
		`s."organization_id"`,
		`s."user_id"`,
		`s."api_key_id"`,
		`s."host_instance_id"`,
		`s."created_at"`,
		`s."updated_at"`,
	).
		From(`"session" s`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *SessionStoreCE) Create(ctx context.Context, m *session.Session) error {
	if _, err := s.builder.
		Insert(`"session"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"user_id"`,
			`"api_key_id"`,
			`"host_instance_id"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.UserID,
			m.APIKeyID,
			m.HostInstanceID,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *SessionStoreCE) Delete(ctx context.Context, m *session.Session) error {
	if _, err := s.builder.
		Delete(`"session"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
