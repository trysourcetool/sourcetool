package group

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Query interface{ isQuery() }

type ByIDQuery struct{ ID uuid.UUID }

func (ByIDQuery) isQuery() {}

func ByID(id uuid.UUID) Query { return ByIDQuery{ID: id} }

type ByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (ByOrganizationIDQuery) isQuery() {}

func ByOrganizationID(organizationID uuid.UUID) Query {
	return ByOrganizationIDQuery{OrganizationID: organizationID}
}

type BySlugQuery struct{ Slug string }

func (BySlugQuery) isQuery() {}

func BySlug(slug string) Query { return BySlugQuery{Slug: slug} }

type BySlugsQuery struct{ Slugs []string }

func (BySlugsQuery) isQuery() {}

func BySlugs(slugs []string) Query { return BySlugsQuery{Slugs: slugs} }

type PageQuery interface{ isPageQuery() }

type PageByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (PageByOrganizationIDQuery) isPageQuery() {}

func PageByOrganizationID(organizationID uuid.UUID) PageQuery {
	return PageByOrganizationIDQuery{OrganizationID: organizationID}
}

type PageByPageIDsQuery struct{ PageIDs []uuid.UUID }

func (PageByPageIDsQuery) isPageQuery() {}

func PageByPageIDs(pageIDs []uuid.UUID) PageQuery {
	return PageByPageIDsQuery{PageIDs: pageIDs}
}

type PageByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (PageByEnvironmentIDQuery) isPageQuery() {}

func PageByEnvironmentID(environmentID uuid.UUID) PageQuery {
	return PageByEnvironmentIDQuery{EnvironmentID: environmentID}
}

type Repository interface {
	Get(context.Context, ...Query) (*Group, error)
	List(context.Context, ...Query) ([]*Group, error)
	Create(context.Context, *Group) error
	Update(context.Context, *Group) error
	Delete(context.Context, *Group) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)

	ListPages(context.Context, ...PageQuery) ([]*GroupPage, error)
	BulkInsertPages(context.Context, []*GroupPage) error
	BulkUpdatePages(context.Context, []*GroupPage) error
	BulkDeletePages(context.Context, []*GroupPage) error
}
