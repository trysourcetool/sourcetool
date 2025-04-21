package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type GroupQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isGroupQuery()
}

type groupByIDQuery struct{ id uuid.UUID }

func (q groupByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."id"`: q.id})
}

func (groupByIDQuery) isGroupQuery() {}

func GroupByID(id uuid.UUID) GroupQuery { return groupByIDQuery{id: id} }

type groupByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q groupByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."organization_id"`: q.organizationID})
}

func (groupByOrganizationIDQuery) isGroupQuery() {}

func GroupByOrganizationID(organizationID uuid.UUID) GroupQuery {
	return groupByOrganizationIDQuery{organizationID: organizationID}
}

type groupBySlugQuery struct{ slug string }

func (q groupBySlugQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: q.slug})
}

func (groupBySlugQuery) isGroupQuery() {}

func GroupBySlug(slug string) GroupQuery { return groupBySlugQuery{slug: slug} }

type groupBySlugsQuery struct{ slugs []string }

func (q groupBySlugsQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: q.slugs})
}

func (groupBySlugsQuery) isGroupQuery() {}

func GroupBySlugs(slugs []string) GroupQuery { return groupBySlugsQuery{slugs: slugs} }

type GroupPageQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isGroupPageQuery()
}

type groupPageByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q groupPageByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = gp."group_id"`).
		Where(sq.Eq{`g."organization_id"`: q.organizationID})
}

func (groupPageByOrganizationIDQuery) isGroupPageQuery() {}

func GroupPageByOrganizationID(organizationID uuid.UUID) GroupPageQuery {
	return groupPageByOrganizationIDQuery{organizationID: organizationID}
}

type groupPageByPageIDsQuery struct{ pageIDs []uuid.UUID }

func (q groupPageByPageIDsQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`gp."page_id"`: q.pageIDs})
}

func (groupPageByPageIDsQuery) isGroupPageQuery() {}

func GroupPageByPageIDs(pageIDs []uuid.UUID) GroupPageQuery {
	return groupPageByPageIDsQuery{pageIDs: pageIDs}
}

type groupPageByEnvironmentIDQuery struct{ environmentID uuid.UUID }

func (q groupPageByEnvironmentIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"page" p ON p."id" = gp."page_id"`).
		Where(sq.Eq{`p."environment_id"`: q.environmentID})
}

func (groupPageByEnvironmentIDQuery) isGroupPageQuery() {}

func GroupPageByEnvironmentID(environmentID uuid.UUID) GroupPageQuery {
	return groupPageByEnvironmentIDQuery{environmentID: environmentID}
}
