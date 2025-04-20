package page

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres/db"
)

type RepositoryCE struct {
	db      db.DB
	builder sq.StatementBuilderType
}

func NewRepositoryCE(db db.DB) *RepositoryCE {
	return &RepositoryCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RepositoryCE) Get(ctx context.Context, queries ...page.Query) (*page.Page, error) {
	query, args, err := r.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := page.Page{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrPageNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) List(ctx context.Context, queries ...page.Query) ([]*page.Page, error) {
	query, args, err := r.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*page.Page, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildQuery(ctx context.Context, queries ...page.Query) (string, []any, error) {
	q := r.builder.Select(
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

	q = r.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) applyQueries(b sq.SelectBuilder, queries ...page.Query) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case page.ByIDQuery:
			b = b.Where(sq.Eq{`p."id"`: q.ID})
		case page.ByOrganizationIDQuery:
			b = b.Where(sq.Eq{`p."organization_id"`: q.OrganizationID})
		case page.ByAPIKeyIDQuery:
			b = b.Where(sq.Eq{`p."api_key_id"`: q.APIKeyID})
		case page.BySessionIDQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."id" = p."api_key_id"`).
				InnerJoin(`"session" s ON s."api_key_id" = ak."id"`).
				Where(sq.Eq{`s."id"`: q.SessionID})
		case page.ByEnvironmentIDQuery:
			b = b.Where(sq.Eq{`p."environment_id"`: q.EnvironmentID})
		case page.LimitQuery:
			b = b.Limit(q.Limit)
		case page.OffsetQuery:
			b = b.Offset(q.Offset)
		case page.OrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}

	return b
}

func (r *RepositoryCE) BulkInsert(ctx context.Context, m []*page.Page) error {
	if len(m) == 0 {
		return nil
	}

	q := r.builder.
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
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) BulkUpdate(ctx context.Context, m []*page.Page) error {
	if len(m) == 0 {
		return nil
	}

	for _, v := range m {
		if _, err := r.builder.
			Update(`"page"`).
			Set(`"name"`, v.Name).
			Set(`"route"`, v.Route).
			Set(`"path"`, v.Path).
			Where(sq.Eq{`"id"`: v.ID}).
			RunWith(r.db).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (r *RepositoryCE) BulkDelete(ctx context.Context, m []*page.Page) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *page.Page, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := r.builder.
		Delete(`"page"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
