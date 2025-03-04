package environment

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
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

func (s *StoreCE) Get(ctx context.Context, opts ...storeopts.EnvironmentOption) (*model.Environment, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.Environment{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrEnvironmentNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) List(ctx context.Context, opts ...storeopts.EnvironmentOption) ([]*model.Environment, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.Environment, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, opts ...storeopts.EnvironmentOption) (string, []any, error) {
	q := s.builder.Select(s.columns()...).
		From(`"environment" e`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *StoreCE) Create(ctx context.Context, m *model.Environment) error {
	if _, err := s.builder.
		Insert(`"environment"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"name"`,
			`"slug"`,
			`"color"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.Name,
			m.Slug,
			m.Color,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) Update(ctx context.Context, m *model.Environment) error {
	if _, err := s.builder.
		Update(`"environment"`).
		Set(`"name"`, m.Name).
		Set(`"slug"`, m.Slug).
		Set(`"color"`, m.Color).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) Delete(ctx context.Context, m *model.Environment) error {
	if _, err := s.builder.
		Delete(`"environment"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) BulkInsert(ctx context.Context, m []*model.Environment) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"environment"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"name"`,
			`"slug"`,
			`"color"`,
		)

	for _, v := range m {
		q = q.Values(
			v.ID,
			v.OrganizationID,
			v.Name,
			v.Slug,
			v.Color,
		)
	}

	if _, err := q.
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) MapByAPIKeyIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*model.Environment, error) {
	cols := append(s.columns(), `ak."id" AS "api_key_id"`)
	query, args, err := s.builder.Select(cols...).
		From(`"environment" e`).
		InnerJoin(`"api_key" ak ON ak."environment_id" = e."id"`).
		Where(sq.Eq{`ak."id"`: ids}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errdefs.ErrDatabase(err)
	}
	defer rows.Close()

	type EnvironmentEmbedded struct {
		*model.Environment
		APIKeyID uuid.UUID `db:"api_key_id"`
	}
	m := make(map[uuid.UUID]*model.Environment)
	for rows.Next() {
		ee := EnvironmentEmbedded{}
		if err := rows.StructScan(&ee); err != nil {
			return nil, errdefs.ErrDatabase(err)
		}

		m[ee.APIKeyID] = ee.Environment
	}

	return m, nil
}

func (s *StoreCE) columns() []string {
	return []string{
		`e."id"`,
		`e."organization_id"`,
		`e."name"`,
		`e."slug"`,
		`e."color"`,
		`e."created_at"`,
		`e."updated_at"`,
	}
}

func (s *StoreCE) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	if _, err := s.Get(ctx, storeopts.EnvironmentByOrganizationID(orgID), storeopts.EnvironmentBySlug(slug)); err != nil {
		if errdefs.IsEnvironmentNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
