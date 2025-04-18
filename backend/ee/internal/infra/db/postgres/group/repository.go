package group

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	groupRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/group"
)

type repositoryEE struct {
	db      db.DB
	builder sq.StatementBuilderType
	*groupRepo.RepositoryCE
}

func NewRepositoryEE(db db.DB) *repositoryEE {
	return &repositoryEE{
		db:           db,
		builder:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		RepositoryCE: groupRepo.NewRepositoryCE(db),
	}
}

func (r *repositoryEE) Get(ctx context.Context, opts ...group.RepositoryOption) (*group.Group, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := group.Group{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *repositoryEE) List(ctx context.Context, opts ...group.RepositoryOption) ([]*group.Group, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*group.Group, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *repositoryEE) buildQuery(ctx context.Context, opts ...group.RepositoryOption) (string, []any, error) {
	q := r.builder.Select(r.columns()...).
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

func (r *repositoryEE) Create(ctx context.Context, m *group.Group) error {
	if _, err := r.builder.
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
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *repositoryEE) Update(ctx context.Context, m *group.Group) error {
	if _, err := r.builder.
		Update(`"group"`).
		Set(`"name"`, m.Name).
		Set(`"slug"`, m.Slug).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *repositoryEE) Delete(ctx context.Context, m *group.Group) error {
	if _, err := r.builder.
		Delete(`"group"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *repositoryEE) columns() []string {
	return []string{
		`g."id"`,
		`g."organization_id"`,
		`g."name"`,
		`g."slug"`,
		`g."created_at"`,
		`g."updated_at"`,
	}
}

func (r *repositoryEE) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := r.Get(ctx, group.ByOrganizationID(orgID), group.BySlug(slug)); err != nil {
		if errdefs.IsGroupNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repositoryEE) ListPages(ctx context.Context, opts ...group.PageRepositoryOption) ([]*group.GroupPage, error) {
	query, args, err := r.buildPageQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*group.GroupPage, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *repositoryEE) buildPageQuery(ctx context.Context, opts ...group.PageRepositoryOption) (string, []any, error) {
	q := r.builder.Select(
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

func (r *repositoryEE) BulkInsertPages(ctx context.Context, pages []*group.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	q := r.builder.
		Insert(`group_page`).
		Columns(
			`"id"`,
			`"group_id"`,
			`"page_id"`,
		)

	for _, p := range pages {
		q = q.Values(p.ID, p.GroupID, p.PageID)
	}

	if _, err := q.RunWith(r.db).ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *repositoryEE) BulkUpdatePages(ctx context.Context, pages []*group.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	for _, p := range pages {
		if _, err := r.builder.
			Update(`group_page`).
			Set(`"group_id"`, p.GroupID).
			Set(`"page_id"`, p.PageID).
			Where(sq.Eq{`"id"`: p.ID}).
			RunWith(r.db).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (r *repositoryEE) BulkDeletePages(ctx context.Context, pages []*group.GroupPage) error {
	if len(pages) == 0 {
		return nil
	}

	ids := lo.Map(pages, func(x *group.GroupPage, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := r.builder.
		Delete(`"group_page"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
