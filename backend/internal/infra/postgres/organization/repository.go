package organization

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
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

func (r *RepositoryCE) Get(ctx context.Context, queries ...organization.Query) (*organization.Organization, error) {
	query, args, err := r.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := organization.Organization{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrOrganizationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) List(ctx context.Context, queries ...organization.Query) ([]*organization.Organization, error) {
	query, args, err := r.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	var orgs []*organization.Organization
	if err := r.db.SelectContext(ctx, &orgs, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return orgs, nil
}

func (r *RepositoryCE) buildQuery(ctx context.Context, queries ...organization.Query) (string, []any, error) {
	q := r.builder.Select(
		`o."id"`,
		`o."subdomain"`,
		`o."created_at"`,
		`o."updated_at"`,
	).
		From(`"organization" o`)

	q = r.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) applyQueries(b sq.SelectBuilder, queries ...organization.Query) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case organization.ByIDQuery:
			b = b.Where(sq.Eq{`o."id"`: q.ID})
		case organization.BySubdomainQuery:
			b = b.Where(sq.Eq{`o."subdomain"`: q.Subdomain})
		case organization.ByUserIDQuery:
			b = b.
				InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
				Where(sq.Eq{`uoa."user_id"`: q.ID})
		}
	}
	return b
}

func (r *RepositoryCE) Create(ctx context.Context, m *organization.Organization) error {
	if _, err := r.builder.
		Insert(`"organization"`).
		Columns(
			`"id"`,
			`"subdomain"`,
		).
		Values(
			m.ID,
			m.Subdomain,
		).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) IsSubdomainExists(ctx context.Context, subdomain string) (bool, error) {
	if _, err := r.Get(ctx, organization.BySubdomain(subdomain)); err != nil {
		if errdefs.IsOrganizationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
