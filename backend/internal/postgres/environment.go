package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.EnvironmentStore = (*environmentStore)(nil)

type environmentStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newEnvironmentStore(db internal.DB) *environmentStore {
	return &environmentStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *environmentStore) Get(ctx context.Context, queries ...database.EnvironmentQuery) (*core.Environment, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Environment{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrEnvironmentNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *environmentStore) List(ctx context.Context, queries ...database.EnvironmentQuery) ([]*core.Environment, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Environment, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *environmentStore) applyQueries(b sq.SelectBuilder, queries ...database.EnvironmentQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.EnvironmentByIDQuery:
			b = b.Where(sq.Eq{`e."id"`: q.ID})
		case database.EnvironmentByOrganizationIDQuery:
			b = b.Where(sq.Eq{`e."organization_id"`: q.OrganizationID})
		case database.EnvironmentBySlugQuery:
			b = b.Where(sq.Eq{`e."slug"`: q.Slug})
		case database.EnvironmentByAPIKeyIDsQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."environment_id" = e."id"`).
				Where(sq.Eq{`ak."id"`: q.APIKeyIDs})
		}
	}
	return b
}

func (s *environmentStore) buildQueryWithColumns(ctx context.Context, extraCols []string, queries ...database.EnvironmentQuery) (string, []any, error) {
	cols := s.columns()
	if len(extraCols) > 0 {
		cols = append(cols, extraCols...)
	}

	q := s.builder.Select(cols...).
		From(`"environment" e`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *environmentStore) buildQuery(ctx context.Context, queries ...database.EnvironmentQuery) (string, []any, error) {
	return s.buildQueryWithColumns(ctx, nil, queries...)
}

func (s *environmentStore) Create(ctx context.Context, m *core.Environment) error {
	if _, err := s.builder.
		Insert(`"environment"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"name"`,
			`"slug"`,
			`"color"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.Name,
			m.Slug,
			m.Color,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *environmentStore) Update(ctx context.Context, m *core.Environment) error {
	if _, err := s.builder.
		Update(`"environment"`).
		Set(`"name"`, m.Name).
		Set(`"slug"`, m.Slug).
		Set(`"color"`, m.Color).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *environmentStore) Delete(ctx context.Context, m *core.Environment) error {
	if _, err := s.builder.
		Delete(`"environment"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *environmentStore) BulkInsert(ctx context.Context, m []*core.Environment) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"environment"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"name"`,
			`"slug"`,
			`"color"`,
		)

	for _, v := range m {
		q = q.Values(
			v.ID,
			v.OrganizationID,
			v.Name,
			v.Slug,
			v.Color,
		)
	}

	if _, err := q.
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *environmentStore) MapByAPIKeyIDs(ctx context.Context, apiKeyIDs []uuid.UUID) (map[uuid.UUID]*core.Environment, error) {
	extraCols := []string{`ak."id" AS "api_key_id"`}
	query, args, err := s.buildQueryWithColumns(ctx, extraCols, database.EnvironmentByAPIKeyIDs(apiKeyIDs))
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errdefs.ErrDatabase(err)
	}
	defer rows.Close()

	type environmentWithAPIKeyID struct {
		*core.Environment
		APIKeyID uuid.UUID `db:"api_key_id"`
	}

	m := make(map[uuid.UUID]*core.Environment)
	for rows.Next() {
		e := environmentWithAPIKeyID{}
		if err := rows.StructScan(&e); err != nil {
			return nil, errdefs.ErrDatabase(err)
		}

		m[e.APIKeyID] = e.Environment
	}

	return m, nil
}

func (s *environmentStore) columns() []string {
	return []string{
		`e."id"`,
		`e."organization_id"`,
		`e."name"`,
		`e."slug"`,
		`e."color"`,
		`e."created_at"`,
		`e."updated_at"`,
	}
}

func (s *environmentStore) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := s.Get(ctx, database.EnvironmentByOrganizationID(orgID), database.EnvironmentBySlug(slug)); err != nil {
		if errdefs.IsEnvironmentNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
