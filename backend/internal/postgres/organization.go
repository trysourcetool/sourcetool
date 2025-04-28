package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.OrganizationStore = (*organizationStore)(nil)

type organizationStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newOrganizationStore(db internal.DB) *organizationStore {
	return &organizationStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *organizationStore) Get(ctx context.Context, queries ...database.OrganizationQuery) (*core.Organization, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Organization{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrOrganizationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *organizationStore) List(ctx context.Context, queries ...database.OrganizationQuery) ([]*core.Organization, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	var orgs []*core.Organization
	if err := s.db.SelectContext(ctx, &orgs, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return orgs, nil
}

func (s *organizationStore) applyQueries(b sq.SelectBuilder, queries ...database.OrganizationQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.OrganizationByIDQuery:
			b = b.Where(sq.Eq{`o."id"`: q.ID})
		case database.OrganizationBySubdomainQuery:
			b = b.Where(sq.Eq{`o."subdomain"`: q.Subdomain})
		case database.OrganizationByUserIDQuery:
			b = b.
				InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
				Where(sq.Eq{`uoa."user_id"`: q.ID})
		}
	}
	return b
}

func (s *organizationStore) buildQuery(ctx context.Context, queries ...database.OrganizationQuery) (string, []any, error) {
	q := s.builder.Select(
		`o."id"`,
		`o."subdomain"`,
		`o."created_at"`,
		`o."updated_at"`,
	).
		From(`"organization" o`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *organizationStore) Create(ctx context.Context, m *core.Organization) error {
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

func (s *organizationStore) IsSubdomainExists(ctx context.Context, subdomain string) (bool, error) {
	if _, err := s.Get(ctx, database.OrganizationBySubdomain(subdomain)); err != nil {
		if errdefs.IsOrganizationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
