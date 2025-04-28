package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.APIKeyStore = (*apiKeyStore)(nil)

type apiKeyStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newAPIKeyStore(db internal.DB) *apiKeyStore {
	return &apiKeyStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *apiKeyStore) Get(ctx context.Context, queries ...database.APIKeyQuery) (*core.APIKey, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.APIKey{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrAPIKeyNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *apiKeyStore) List(ctx context.Context, queries ...database.APIKeyQuery) ([]*core.APIKey, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.APIKey, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *apiKeyStore) applyQueries(b sq.SelectBuilder, queries ...database.APIKeyQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.APIKeyByIDQuery:
			b = b.Where(sq.Eq{`ak."id"`: q.ID})
		case database.APIKeyByOrganizationIDQuery:
			b = b.Where(sq.Eq{`ak."organization_id"`: q.OrganizationID})
		case database.APIKeyByEnvironmentIDQuery:
			b = b.Where(sq.Eq{`ak."environment_id"`: q.EnvironmentID})
		case database.APIKeyByEnvironmentIDsQuery:
			b = b.Where(sq.Eq{`ak."environment_id"`: q.EnvironmentIDs})
		case database.APIKeyByUserIDQuery:
			b = b.Where(sq.Eq{`ak."user_id"`: q.UserID})
		case database.APIKeyByKeyQuery:
			b = b.Where(sq.Eq{`ak."key"`: q.Key})
		}
	}
	return b
}

func (s *apiKeyStore) buildQuery(ctx context.Context, queries ...database.APIKeyQuery) (string, []any, error) {
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

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *apiKeyStore) Create(ctx context.Context, m *core.APIKey) error {
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

func (s *apiKeyStore) Update(ctx context.Context, m *core.APIKey) error {
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

func (s *apiKeyStore) Delete(ctx context.Context, m *core.APIKey) error {
	if _, err := s.builder.
		Delete(`"api_key"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
