package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type OrganizationQuery interface{ isOrganizationQuery() }

type OrganizationByIDQuery struct{ ID uuid.UUID }

func (OrganizationByIDQuery) isOrganizationQuery() {}

func OrganizationByID(id uuid.UUID) OrganizationQuery { return OrganizationByIDQuery{ID: id} }

type OrganizationBySubdomainQuery struct{ Subdomain string }

func (OrganizationBySubdomainQuery) isOrganizationQuery() {}

func OrganizationBySubdomain(subdomain string) OrganizationQuery {
	return OrganizationBySubdomainQuery{Subdomain: subdomain}
}

type OrganizationByUserIDQuery struct{ ID uuid.UUID }

func (OrganizationByUserIDQuery) isOrganizationQuery() {}

func OrganizationByUserID(id uuid.UUID) OrganizationQuery { return OrganizationByUserIDQuery{ID: id} }

func applyOrganizationQueries(b sq.SelectBuilder, queries ...OrganizationQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case OrganizationByIDQuery:
			b = b.Where(sq.Eq{`o."id"`: q.ID})
		case OrganizationBySubdomainQuery:
			b = b.Where(sq.Eq{`o."subdomain"`: q.Subdomain})
		case OrganizationByUserIDQuery:
			b = b.
				InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
				Where(sq.Eq{`uoa."user_id"`: q.ID})
		}
	}
	return b
}
