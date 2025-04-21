package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type EnvironmentQuery interface{ isEnvironmentQuery() }

type EnvironmentByIDQuery struct{ ID uuid.UUID }

func (EnvironmentByIDQuery) isEnvironmentQuery() {}

func EnvironmentByID(id uuid.UUID) EnvironmentQuery { return EnvironmentByIDQuery{ID: id} }

type EnvironmentByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (EnvironmentByOrganizationIDQuery) isEnvironmentQuery() {}

func EnvironmentByOrganizationID(organizationID uuid.UUID) EnvironmentQuery {
	return EnvironmentByOrganizationIDQuery{OrganizationID: organizationID}
}

type EnvironmentBySlugQuery struct{ Slug string }

func (EnvironmentBySlugQuery) isEnvironmentQuery() {}

func EnvironmentBySlug(slug string) EnvironmentQuery { return EnvironmentBySlugQuery{Slug: slug} }

func applyEnvironmentQueries(b sq.SelectBuilder, queries ...EnvironmentQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case EnvironmentByIDQuery:
			b = b.Where(sq.Eq{`e."id"`: q.ID})
		case EnvironmentByOrganizationIDQuery:
			b = b.Where(sq.Eq{`e."organization_id"`: q.OrganizationID})
		case EnvironmentBySlugQuery:
			b = b.Where(sq.Eq{`e."slug"`: q.Slug})
		}
	}
	return b
}
