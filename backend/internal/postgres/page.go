package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (db *DB) GetPage(ctx context.Context, queries ...PageQuery) (*core.Page, error) {
	query, args, err := db.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Page{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrPageNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListPages(ctx context.Context, queries ...PageQuery) ([]*core.Page, error) {
	query, args, err := db.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Page, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildQuery(ctx context.Context, queries ...PageQuery) (string, []any, error) {
	q := db.builder.Select(
		`p."id"`,
		`p."organization_id"`,
		`p."environment_id"`,
		`p."api_key_id"`,
		`p."name"`,
		`p."route"`,
		`p."path"`,
		`p."created_at"`,
		`p."updated_at"`,
	).
		From(`"page" p`)

	q = applyPageQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) BulkInsertPages(ctx context.Context, tx *sqlx.Tx, m []*core.Page) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(m) == 0 {
		return nil
	}

	q := db.builder.
		Insert(`"page"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"environment_id"`,
			`"api_key_id"`,
			`"name"`,
			`"route"`,
			`"path"`,
		)

	for _, v := range m {
		q = q.Values(
			v.ID,
			v.OrganizationID,
			v.EnvironmentID,
			v.APIKeyID,
			v.Name,
			v.Route,
			v.Path,
		)
	}

	if _, err := q.
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) BulkUpdatePages(ctx context.Context, tx *sqlx.Tx, m []*core.Page) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(m) == 0 {
		return nil
	}

	for _, v := range m {
		if _, err := db.builder.
			Update(`"page"`).
			Set(`"name"`, v.Name).
			Set(`"route"`, v.Route).
			Set(`"path"`, v.Path).
			Where(sq.Eq{`"id"`: v.ID}).
			RunWith(runner).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (db *DB) BulkDeletePages(ctx context.Context, tx *sqlx.Tx, m []*core.Page) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *core.Page, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := db.builder.
		Delete(`"page"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
