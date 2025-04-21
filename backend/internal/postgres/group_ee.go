//go:build ee
// +build ee

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

func (db *DB) GetGroup(ctx context.Context, queries ...GroupQuery) (*core.Group, error) {
	query, args, err := db.buildGroupQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Group{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListGroups(ctx context.Context, queries ...GroupQuery) ([]*core.Group, error) {
	query, args, err := db.buildGroupQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Group, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildGroupQuery(ctx context.Context, queries ...GroupQuery) (string, []any, error) {
	q := db.builder.Select(db.columns()...).
		From(`"group" g`)

	for _, query := range queries {
		q = query.apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) CreateGroup(ctx context.Context, tx *sqlx.Tx, m *core.Group) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Insert(`"group"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"name"`,
			`"slug"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.Name,
			m.Slug,
		).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) UpdateGroup(ctx context.Context, tx *sqlx.Tx, m *core.Group) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Update(`"group"`).
		Set(`"name"`, m.Name).
		Set(`"slug"`, m.Slug).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) DeleteGroup(ctx context.Context, tx *sqlx.Tx, m *core.Group) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Delete(`"group"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) columns() []string {
	return []string{
		`g."id"`,
		`g."organization_id"`,
		`g."name"`,
		`g."slug"`,
		`g."created_at"`,
		`g."updated_at"`,
	}
}

func (db *DB) IsGroupSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := db.GetGroup(ctx, GroupByOrganizationID(orgID), GroupBySlug(slug)); err != nil {
		if errdefs.IsGroupNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *DB) ListGroupPages(ctx context.Context, queries ...GroupPageQuery) ([]*core.GroupPage, error) {
	query, args, err := db.buildGroupPageQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.GroupPage, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildGroupPageQuery(ctx context.Context, queries ...GroupPageQuery) (string, []any, error) {
	q := db.builder.Select(
		`gp."id"`,
		`gp."group_id"`,
		`gp."page_id"`,
		`gp."created_at"`,
		`gp."updated_at"`,
	).
		From(`"group_page" gp`)

	for _, query := range queries {
		q = query.apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) BulkInsertGroupPages(ctx context.Context, tx *sqlx.Tx, pages []*core.GroupPage) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(pages) == 0 {
		return nil
	}

	q := db.builder.
		Insert(`group_page`).
		Columns(
			`"id"`,
			`"group_id"`,
			`"page_id"`,
		)

	for _, p := range pages {
		q = q.Values(p.ID, p.GroupID, p.PageID)
	}

	if _, err := q.RunWith(runner).ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) BulkUpdateGroupPages(ctx context.Context, tx *sqlx.Tx, pages []*core.GroupPage) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(pages) == 0 {
		return nil
	}

	for _, p := range pages {
		if _, err := db.builder.
			Update(`group_page`).
			Set(`"group_id"`, p.GroupID).
			Set(`"page_id"`, p.PageID).
			Where(sq.Eq{`"id"`: p.ID}).
			RunWith(runner).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (db *DB) BulkDeleteGroupPages(ctx context.Context, tx *sqlx.Tx, pages []*core.GroupPage) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(pages) == 0 {
		return nil
	}

	ids := lo.Map(pages, func(x *core.GroupPage, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := db.builder.
		Delete(`"group_page"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
