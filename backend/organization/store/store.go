package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization"
)

type OrganizationStoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewOrganizationStoreCE(db infra.DB) *OrganizationStoreCE {
	return &OrganizationStoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *OrganizationStoreCE) Get(ctx context.Context, opts ...organization.StoreOption) (*organization.Organization, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := organization.Organization{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrOrganizationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *OrganizationStoreCE) List(ctx context.Context, opts ...organization.StoreOption) ([]*organization.Organization, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var orgs []*organization.Organization
	if err := s.db.SelectContext(ctx, &orgs, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return orgs, nil
}

func (s *OrganizationStoreCE) buildQuery(ctx context.Context, opts ...organization.StoreOption) (string, []any, error) {
	q := s.builder.Select(
		`o."id"`,
		`o."subdomain"`,
		`o."created_at"`,
		`o."updated_at"`,
	).
		From(`"organization" o`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *OrganizationStoreCE) Create(ctx context.Context, m *organization.Organization) error {
	if _, err := s.builder.
		Insert(`"organization"`).
		Columns(
			`"id"`,
			`"subdomain"`,
		).
		Values(
			m.ID,
			m.Subdomain,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *OrganizationStoreCE) IsSubdomainExists(ctx context.Context, subdomain string) (bool, error) {
	if _, err := s.Get(ctx, organization.BySubdomain(subdomain)); err != nil {
		if errdefs.IsOrganizationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
