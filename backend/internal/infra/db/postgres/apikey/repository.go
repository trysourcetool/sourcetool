package apikey

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
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

func (r *RepositoryCE) Get(ctx context.Context, opts ...apikey.RepositoryOption) (*apikey.APIKey, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := apikey.APIKey{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrAPIKeyNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) List(ctx context.Context, opts ...apikey.RepositoryOption) ([]*apikey.APIKey, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*apikey.APIKey, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildQuery(ctx context.Context, opts ...apikey.RepositoryOption) (string, []any, error) {
	q := r.builder.Select(
		`ak."id"`,
		`ak."organization_id"`,
		`ak."environment_id"`,
		`ak."user_id"`,
		`ak."name"`,
		`ak."key"`,
		`ak."created_at"`,
		`ak."updated_at"`,
	).
		From(`"api_key" ak`)

	for _, opt := range opts {
		q = opt.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) Create(ctx context.Context, m *apikey.APIKey) error {
	if _, err := r.builder.
		Insert(`"api_key"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"environment_id"`,
			`"user_id"`,
			`"name"`,
			`"key"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.EnvironmentID,
			m.UserID,
			m.Name,
			m.Key,
		).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) Update(ctx context.Context, m *apikey.APIKey) error {
	if _, err := r.builder.
		Update(`"api_key"`).
		Set(`"user_id"`, m.UserID).
		Set(`"name"`, m.Name).
		Set(`"key"`, m.Key).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) Delete(ctx context.Context, m *apikey.APIKey) error {
	if _, err := r.builder.
		Delete(`"api_key"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
