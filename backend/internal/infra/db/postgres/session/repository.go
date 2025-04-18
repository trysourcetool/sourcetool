package session

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
)

type RepositoryCE struct {
	db      db.DB
	builder sq.StatementBuilderType
}

func NewRepositoryCE(db db.DB) *RepositoryCE {
	return &RepositoryCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RepositoryCE) Get(ctx context.Context, opts ...session.RepositoryOption) (*session.Session, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := session.Session{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrSessionNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) buildQuery(ctx context.Context, opts ...session.RepositoryOption) (string, []any, error) {
	q := r.builder.Select(
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

func (r *RepositoryCE) Create(ctx context.Context, m *session.Session) error {
	if _, err := r.builder.
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
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) Delete(ctx context.Context, m *session.Session) error {
	if _, err := r.builder.
		Delete(`"session"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
