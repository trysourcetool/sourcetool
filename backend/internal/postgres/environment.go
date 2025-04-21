package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (db *DB) GetEnvironment(ctx context.Context, queries ...EnvironmentQuery) (*core.Environment, error) {
	query, args, err := db.buildEnvironmentQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Environment{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrEnvironmentNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListEnvironments(ctx context.Context, queries ...EnvironmentQuery) ([]*core.Environment, error) {
	query, args, err := db.buildEnvironmentQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Environment, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildEnvironmentQuery(ctx context.Context, queries ...EnvironmentQuery) (string, []any, error) {
	q := db.builder.Select(environmentColumns()...).
		From(`"environment" e`)

	q = applyEnvironmentQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) CreateEnvironment(ctx context.Context, tx *sqlx.Tx, m *core.Environment) error {
	if _, err := db.builder.
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
		RunWith(tx).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) UpdateEnvironment(ctx context.Context, tx *sqlx.Tx, m *core.Environment) error {
	if _, err := db.builder.
		Update(`"environment"`).
		Set(`"name"`, m.Name).
		Set(`"slug"`, m.Slug).
		Set(`"color"`, m.Color).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(tx).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) DeleteEnvironment(ctx context.Context, tx *sqlx.Tx, m *core.Environment) error {
	if _, err := db.builder.
		Delete(`"environment"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(tx).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) BulkInsertEnvironments(ctx context.Context, tx *sqlx.Tx, m []*core.Environment) error {
	if len(m) == 0 {
		return nil
	}

	q := db.builder.
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
		RunWith(tx).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) MapEnvironmentsByAPIKeyIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*core.Environment, error) {
	cols := append(environmentColumns(), `ak."id" AS "api_key_id"`)
	query, args, err := db.builder.Select(cols...).
		From(`"environment" e`).
		InnerJoin(`"api_key" ak ON ak."environment_id" = e."id"`).
		Where(sq.Eq{`ak."id"`: ids}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errdefs.ErrDatabase(err)
	}
	defer rows.Close()

	type EnvironmentEmbedded struct {
		*core.Environment
		APIKeyID uuid.UUID `db:"api_key_id"`
	}
	m := make(map[uuid.UUID]*core.Environment)
	for rows.Next() {
		ee := EnvironmentEmbedded{}
		if err := rows.StructScan(&ee); err != nil {
			return nil, errdefs.ErrDatabase(err)
		}

		m[ee.APIKeyID] = ee.Environment
	}

	return m, nil
}

func environmentColumns() []string {
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

func (db *DB) IsEnvironmentSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := db.GetEnvironment(ctx, EnvironmentByOrganizationID(orgID), EnvironmentBySlug(slug)); err != nil {
		if errdefs.IsEnvironmentNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
