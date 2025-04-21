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

var _ database.PageStore = (*pageStore)(nil)

type pageStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newPageStore(db internal.DB) *pageStore {
	return &pageStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *pageStore) Get(ctx context.Context, queries ...database.PageQuery) (*core.Page, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Page{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrPageNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *pageStore) List(ctx context.Context, queries ...database.PageQuery) ([]*core.Page, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Page, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *pageStore) applyQueries(b sq.SelectBuilder, queries ...database.PageQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.PageByIDQuery:
			b = b.Where(sq.Eq{`p."id"`: q.ID})
		case database.PageByOrganizationIDQuery:
			b = b.Where(sq.Eq{`p."organization_id"`: q.OrganizationID})
		case database.PageByAPIKeyIDQuery:
			b = b.Where(sq.Eq{`p."api_key_id"`: q.APIKeyID})
		case database.PageBySessionIDQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."id" = p."api_key_id"`).
				InnerJoin(`"environment" e ON e."id" = ak."environment_id"`).
				InnerJoin(`"session" s ON s."environment_id" = e."id"`).
				Where(sq.Eq{`s."id"`: q.SessionID})
		case database.PageByEnvironmentIDQuery:
			b = b.Where(sq.Eq{`p."environment_id"`: q.EnvironmentID})
		case database.PageLimitQuery:
			b = b.Limit(q.Limit)
		case database.PageOffsetQuery:
			b = b.Offset(q.Offset)
		case database.PageOrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}

	return b
}

func (s *pageStore) buildQuery(ctx context.Context, queries ...database.PageQuery) (string, []any, error) {
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

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *pageStore) BulkInsert(ctx context.Context, m []*core.Page) error {
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

func (s *pageStore) BulkUpdate(ctx context.Context, m []*core.Page) error {
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

func (s *pageStore) BulkDelete(ctx context.Context, m []*core.Page) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *core.Page, _ int) uuid.UUID {
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
