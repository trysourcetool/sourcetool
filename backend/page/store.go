package page

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type StoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewStoreCE(db infra.DB) *StoreCE {
	return &StoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *StoreCE) Get(ctx context.Context, conditions ...any) (*model.Page, error) {
	query, args, err := s.buildQuery(ctx, conditions...)
	if err != nil {
		return nil, err
	}

	m := model.Page{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrPageNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) List(ctx context.Context, conditions ...any) ([]*model.Page, error) {
	query, args, err := s.buildQuery(ctx, conditions...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.Page, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, conditions ...any) (string, []any, error) {
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

	opts, err := s.toSelectOptions(ctx, conditions...)
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	for _, o := range opts {
		q = o(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *StoreCE) toSelectOptions(ctx context.Context, conditions ...any) ([]infra.SelectOption, error) {
	options := make([]infra.SelectOption, len(conditions))
	for i, c := range conditions {
		switch v := c.(type) {
		case model.PageByID:
			options[i] = s.byID(v)
		case model.PageByOrganizationID:
			options[i] = s.byOrganizationID(v)
		case model.PageByAPIKeyID:
			options[i] = s.byAPIKeyID(v)
		case model.PageBySessionID:
			options[i] = s.bySessionID(v)
		case infra.Limit:
			options[i] = s.limit(v)
		case infra.Offset:
			options[i] = s.offset(v)
		case infra.OrderBy:
			options[i] = s.orderBy(v)
		default:
			return nil, errdefs.ErrDatabase(errors.New("unsupported condition"))
		}
	}

	return options, nil
}

func (s *StoreCE) byID(in model.PageByID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`p."id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) byOrganizationID(in model.PageByOrganizationID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`p."organization_id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) byAPIKeyID(in model.PageByAPIKeyID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`p."api_key_id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) bySessionID(in model.PageBySessionID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.
			InnerJoin(`"session" s ON s."page_id" = p."id"`).
			Where(sq.Eq{`s."id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) limit(in infra.Limit) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Limit(uint64(in))
	}
}

func (s *StoreCE) offset(in infra.Offset) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Offset(uint64(in))
	}
}

func (s *StoreCE) orderBy(in infra.OrderBy) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.OrderBy(string(in))
	}
}

func (s *StoreCE) BulkInsert(ctx context.Context, m []*model.Page) error {
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

func (s *StoreCE) BulkUpdate(ctx context.Context, m []*model.Page) error {
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

func (s *StoreCE) BulkDelete(ctx context.Context, m []*model.Page) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *model.Page, _ int) uuid.UUID {
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
