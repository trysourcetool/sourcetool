package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/page"
)

type PageStoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewPageStoreCE(db infra.DB) *PageStoreCE {
	return &PageStoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *PageStoreCE) Get(ctx context.Context, opts ...page.StoreOption) (*page.Page, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := page.Page{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrPageNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *PageStoreCE) List(ctx context.Context, opts ...page.StoreOption) ([]*page.Page, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*page.Page, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *PageStoreCE) buildQuery(ctx context.Context, opts ...page.StoreOption) (string, []any, error) {
	q := s.builder.Select(
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

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *PageStoreCE) BulkInsert(ctx context.Context, m []*page.Page) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
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
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *PageStoreCE) BulkUpdate(ctx context.Context, m []*page.Page) error {
	if len(m) == 0 {
		return nil
	}

	for _, v := range m {
		if _, err := s.builder.
			Update(`"page"`).
			Set(`"name"`, v.Name).
			Set(`"route"`, v.Route).
			Set(`"path"`, v.Path).
			Where(sq.Eq{`"id"`: v.ID}).
			RunWith(s.db).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (s *PageStoreCE) BulkDelete(ctx context.Context, m []*page.Page) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *page.Page, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := s.builder.
		Delete(`"page"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
