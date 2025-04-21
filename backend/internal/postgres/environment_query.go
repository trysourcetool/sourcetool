package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type EnvironmentQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isEnvironmentQuery()
}

type environmentByIDQuery struct{ id uuid.UUID }

func (q environmentByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."id"`: q.id})
}

func (environmentByIDQuery) isEnvironmentQuery() {}

func EnvironmentByID(id uuid.UUID) EnvironmentQuery { return environmentByIDQuery{id: id} }

type environmentByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q environmentByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."organization_id"`: q.organizationID})
}

func (environmentByOrganizationIDQuery) isEnvironmentQuery() {}

func EnvironmentByOrganizationID(organizationID uuid.UUID) EnvironmentQuery {
	return environmentByOrganizationIDQuery{organizationID: organizationID}
}

type environmentBySlugQuery struct{ slug string }

func (q environmentBySlugQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."slug"`: q.slug})
}

func (environmentBySlugQuery) isEnvironmentQuery() {}

func EnvironmentBySlug(slug string) EnvironmentQuery { return environmentBySlugQuery{slug: slug} }
