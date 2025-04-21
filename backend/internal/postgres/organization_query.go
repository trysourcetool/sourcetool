package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type OrganizationQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isOrganizationQuery()
}

type organizationByIDQuery struct{ id uuid.UUID }

func (q organizationByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`o."id"`: q.id})
}

func (organizationByIDQuery) isOrganizationQuery() {}

func OrganizationByID(id uuid.UUID) OrganizationQuery { return organizationByIDQuery{id: id} }

type organizationBySubdomainQuery struct{ subdomain string }

func (q organizationBySubdomainQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`o."subdomain"`: q.subdomain})
}

func (organizationBySubdomainQuery) isOrganizationQuery() {}

func OrganizationBySubdomain(subdomain string) OrganizationQuery {
	return organizationBySubdomainQuery{subdomain: subdomain}
}

type organizationByUserIDQuery struct{ id uuid.UUID }

func (q organizationByUserIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
		Where(sq.Eq{`uoa."user_id"`: q.id})
}

func (organizationByUserIDQuery) isOrganizationQuery() {}

func OrganizationByUserID(id uuid.UUID) OrganizationQuery { return organizationByUserIDQuery{id: id} }
