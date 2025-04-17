package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type APIKeyStoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewAPIKeyStoreCE(db infra.DB) *APIKeyStoreCE {
	return &APIKeyStoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *APIKeyStoreCE) Get(ctx context.Context, opts ...storeopts.APIKeyOption) (*model.APIKey, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.APIKey{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrAPIKeyNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *APIKeyStoreCE) List(ctx context.Context, opts ...storeopts.APIKeyOption) ([]*model.APIKey, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.APIKey, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *APIKeyStoreCE) buildQuery(ctx context.Context, opts ...storeopts.APIKeyOption) (string, []any, error) {
	q := s.builder.Select(
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

func (s *APIKeyStoreCE) Create(ctx context.Context, m *model.APIKey) error {
	if _, err := s.builder.
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
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *APIKeyStoreCE) Update(ctx context.Context, m *model.APIKey) error {
	if _, err := s.builder.
		Update(`"api_key"`).
		Set(`"user_id"`, m.UserID).
		Set(`"name"`, m.Name).
		Set(`"key"`, m.Key).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *APIKeyStoreCE) Delete(ctx context.Context, m *model.APIKey) error {
	if _, err := s.builder.
		Delete(`"api_key"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
