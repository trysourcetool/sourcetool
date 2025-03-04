package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type OrganizationOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isOrganizationOption()
}

func OrganizationByID(id uuid.UUID) OrganizationOption {
	return organizationByIDOption{id: id}
}

type organizationByIDOption struct {
	id uuid.UUID
}

func (o organizationByIDOption) isOrganizationOption() {}

func (o organizationByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`o."id"`: o.id})
}

func OrganizationBySubdomain(subdomain string) OrganizationOption {
	return organizationBySubdomainOption{subdomain: subdomain}
}

type organizationBySubdomainOption struct {
	subdomain string
}

func (o organizationBySubdomainOption) isOrganizationOption() {}

func (o organizationBySubdomainOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`o."subdomain"`: o.subdomain})
}

func OrganizationByUserID(id uuid.UUID) OrganizationOption {
	return organizationByUserIDOption{id: id}
}

type organizationByUserIDOption struct {
	id uuid.UUID
}

func (o organizationByUserIDOption) isOrganizationOption() {}

func (o organizationByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
		Where(sq.Eq{`uoa."user_id"`: o.id})
}
