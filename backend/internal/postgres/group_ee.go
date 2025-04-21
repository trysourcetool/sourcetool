//go:build ee
// +build ee

package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.GroupStore = (*groupStore)(nil)

type groupStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newGroupStore(db internal.DB) *groupStore {
	return &groupStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *groupStore) Get(ctx context.Context, queries ...database.GroupQuery) (*core.Group, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Group{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *groupStore) List(ctx context.Context, queries ...database.GroupQuery) ([]*core.Group, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Group, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *groupStore) applyQueries(b sq.SelectBuilder, queries ...database.GroupQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.GroupByIDQuery:
			b = b.Where(sq.Eq{`g."id"`: q.ID})
		case database.GroupByOrganizationIDQuery:
			b = b.Where(sq.Eq{`g."organization_id"`: q.OrganizationID})
		case database.GroupBySlugQuery:
			b = b.Where(sq.Eq{`g."slug"`: q.Slug})
		case database.GroupBySlugsQuery:
			b = b.Where(sq.Eq{`g."slug"`: q.Slugs})
		}
	}

	return b
}

func (s *groupStore) buildQuery(ctx context.Context, queries ...database.GroupQuery) (string, []any, error) {
	q := s.builder.Select(s.columns()...).
		From(`"group" g`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *groupStore) Create(ctx context.Context, m *core.Group) error {
	if _, err := s.builder.
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
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *groupStore) Update(ctx context.Context, m *core.Group) error {
	if _, err := s.builder.
		Update(`"group"`).
		Set(`"name"`, m.Name).
		Set(`"slug"`, m.Slug).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *groupStore) Delete(ctx context.Context, m *core.Group) error {
	if _, err := s.builder.
		Delete(`"group"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *groupStore) columns() []string {
	return []string{
		`g."id"`,
		`g."organization_id"`,
		`g."name"`,
		`g."slug"`,
		`g."created_at"`,
		`g."updated_at"`,
	}
}

func (s *groupStore) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := s.Get(ctx, database.GroupByOrganizationID(orgID), database.GroupBySlug(slug)); err != nil {
		if errdefs.IsGroupNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *groupStore) ListPages(ctx context.Context, queries ...database.GroupPageQuery) ([]*core.GroupPage, error) {
	query, args, err := s.buildPageQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.GroupPage, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *groupStore) applyPageQueries(b sq.SelectBuilder, queries ...database.GroupPageQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.GroupPageByOrganizationIDQuery:
			b = b.
				InnerJoin(`"group" g ON g."id" = gp."group_id"`).
				Where(sq.Eq{`g."organization_id"`: q.OrganizationID})
		case database.GroupPageByPageIDsQuery:
			b = b.Where(sq.Eq{`gp."page_id"`: q.PageIDs})
		case database.GroupPageByEnvironmentIDQuery:
			b = b.
				InnerJoin(`"page" p ON p."id" = gp."page_id"`).
				Where(sq.Eq{`p."environment_id"`: q.EnvironmentID})
		}
	}

	return b
}

func (s *groupStore) buildPageQuery(ctx context.Context, queries ...database.GroupPageQuery) (string, []any, error) {
	q := s.builder.Select(
		`gp."id"`,
		`gp."group_id"`,
		`gp."page_id"`,
		`gp."created_at"`,
		`gp."updated_at"`,
	).
		From(`"group_page" gp`)

	q = s.applyPageQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *groupStore) BulkInsertPages(ctx context.Context, pages []*core.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"group_page"`).
		Columns(
			`"id"`,
			`"group_id"`,
			`"page_id"`,
		)

	for _, p := range pages {
		q = q.Values(p.ID, p.GroupID, p.PageID)
	}

	if _, err := q.RunWith(s.db).ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *groupStore) BulkUpdatePages(ctx context.Context, pages []*core.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	for _, p := range pages {
		if _, err := s.builder.
			Update(`"group_page"`).
			Set(`"group_id"`, p.GroupID).
			Set(`"page_id"`, p.PageID).
			Where(sq.Eq{`"id"`: p.ID}).
			RunWith(s.db).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (s *groupStore) BulkDeletePages(ctx context.Context, pages []*core.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	ids := lo.Map(pages, func(x *core.GroupPage, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := s.builder.
		Delete(`"group_page"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
