package postgres

import (
	"github.com/gofrs/uuid/v5"
)

type GroupQuery interface{ isGroupQuery() }

type GroupByIDQuery struct{ ID uuid.UUID }

func (GroupByIDQuery) isGroupQuery() {}

func GroupByID(id uuid.UUID) GroupQuery { return GroupByIDQuery{ID: id} }

type GroupByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (GroupByOrganizationIDQuery) isGroupQuery() {}

func GroupByOrganizationID(organizationID uuid.UUID) GroupQuery {
	return GroupByOrganizationIDQuery{OrganizationID: organizationID}
}

type GroupBySlugQuery struct{ Slug string }

func (GroupBySlugQuery) isGroupQuery() {}

func GroupBySlug(slug string) GroupQuery { return GroupBySlugQuery{Slug: slug} }

type GroupBySlugsQuery struct{ Slugs []string }

func (GroupBySlugsQuery) isGroupQuery() {}

func GroupBySlugs(slugs []string) GroupQuery { return GroupBySlugsQuery{Slugs: slugs} }

type GroupPageQuery interface{ isGroupPageQuery() }

type GroupPageByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (GroupPageByOrganizationIDQuery) isGroupPageQuery() {}

func GroupPageByOrganizationID(organizationID uuid.UUID) GroupPageQuery {
	return GroupPageByOrganizationIDQuery{OrganizationID: organizationID}
}

type GroupPageByPageIDsQuery struct{ PageIDs []uuid.UUID }

func (GroupPageByPageIDsQuery) isGroupPageQuery() {}

func GroupPageByPageIDs(pageIDs []uuid.UUID) GroupPageQuery {
	return GroupPageByPageIDsQuery{PageIDs: pageIDs}
}

type GroupPageByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (GroupPageByEnvironmentIDQuery) isGroupPageQuery() {}

func GroupPageByEnvironmentID(environmentID uuid.UUID) GroupPageQuery {
	return GroupPageByEnvironmentIDQuery{EnvironmentID: environmentID}
}
