package group

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/group/store"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type groupStoreEE struct {
	db      infra.DB
	builder sq.StatementBuilderType
	*store.GroupStoreCE
}

func NewGroupStoreEE(db infra.DB) *groupStoreEE {
	return &groupStoreEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		GroupStoreCE: store.NewGroupStoreCE(db),
	}
}

func (s *groupStoreEE) Get(ctx context.Context, opts ...group.StoreOption) (*group.Group, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := group.Group{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *groupStoreEE) List(ctx context.Context, opts ...group.StoreOption) ([]*group.Group, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*group.Group, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *groupStoreEE) buildQuery(ctx context.Context, opts ...group.StoreOption) (string, []any, error) {
	q := s.builder.Select(s.columns()...).
		From(`"group" g`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *groupStoreEE) Create(ctx context.Context, m *group.Group) error {
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

func (s *groupStoreEE) Update(ctx context.Context, m *group.Group) error {
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

func (s *groupStoreEE) Delete(ctx context.Context, m *group.Group) error {
	if _, err := s.builder.
		Delete(`"group"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *groupStoreEE) columns() []string {
	return []string{
		`g."id"`,
		`g."organization_id"`,
		`g."name"`,
		`g."slug"`,
		`g."created_at"`,
		`g."updated_at"`,
	}
}

func (s *groupStoreEE) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := s.Get(ctx, group.ByOrganizationID(orgID), group.BySlug(slug)); err != nil {
		if errdefs.IsGroupNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *groupStoreEE) ListPages(ctx context.Context, opts ...group.PageStoreOption) ([]*group.GroupPage, error) {
	query, args, err := s.buildPageQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*group.GroupPage, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *groupStoreEE) buildPageQuery(ctx context.Context, opts ...group.PageStoreOption) (string, []any, error) {
	q := s.builder.Select(
		`gp."id"`,
		`gp."group_id"`,
		`gp."page_id"`,
		`gp."created_at"`,
		`gp."updated_at"`,
	).
		From(`"group_page" gp`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *groupStoreEE) BulkInsertPages(ctx context.Context, pages []*group.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`group_page`).
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

func (s *groupStoreEE) BulkUpdatePages(ctx context.Context, pages []*group.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	for _, p := range pages {
		if _, err := s.builder.
			Update(`group_page`).
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

func (s *groupStoreEE) BulkDeletePages(ctx context.Context, pages []*group.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	ids := lo.Map(pages, func(x *group.GroupPage, _ int) uuid.UUID {
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
