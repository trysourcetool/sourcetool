package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
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

type GroupStore interface {
	Get(ctx context.Context, queries ...GroupQuery) (*core.Group, error)
	List(ctx context.Context, queries ...GroupQuery) ([]*core.Group, error)
	Create(ctx context.Context, m *core.Group) error
	Update(ctx context.Context, m *core.Group) error
	Delete(ctx context.Context, m *core.Group) error
	IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error)

	ListPages(ctx context.Context, queries ...GroupPageQuery) ([]*core.GroupPage, error)
	BulkInsertPages(ctx context.Context, pages []*core.GroupPage) error
	BulkUpdatePages(ctx context.Context, pages []*core.GroupPage) error
	BulkDeletePages(ctx context.Context, pages []*core.GroupPage) error
}
