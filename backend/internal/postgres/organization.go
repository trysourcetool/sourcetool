package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (db *DB) GetOrganization(ctx context.Context, queries ...OrganizationQuery) (*core.Organization, error) {
	query, args, err := db.buildOrganizationQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Organization{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrOrganizationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListOrganizations(ctx context.Context, queries ...OrganizationQuery) ([]*core.Organization, error) {
	query, args, err := db.buildOrganizationQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	var orgs []*core.Organization
	if err := db.db.SelectContext(ctx, &orgs, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return orgs, nil
}

func (db *DB) buildOrganizationQuery(ctx context.Context, queries ...OrganizationQuery) (string, []any, error) {
	q := db.builder.Select(
		`o."id"`,
		`o."subdomain"`,
		`o."created_at"`,
		`o."updated_at"`,
	).
		From(`"organization" o`)

	for _, query := range queries {
		q = query.apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) CreateOrganization(ctx context.Context, tx *sqlx.Tx, m *core.Organization) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Insert(`"organization"`).
		Columns(
			`"id"`,
			`"subdomain"`,
		).
		Values(
			m.ID,
			m.Subdomain,
		).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) IsOrganizationSubdomainExists(ctx context.Context, subdomain string) (bool, error) {
	if _, err := db.GetOrganization(ctx, OrganizationBySubdomain(subdomain)); err != nil {
		if errdefs.IsOrganizationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
