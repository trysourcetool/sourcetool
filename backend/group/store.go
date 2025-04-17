package group

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type StoreOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isStoreOption()
}

func ByID(id uuid.UUID) StoreOption {
	return byIDOption{id: id}
}

type byIDOption struct {
	id uuid.UUID
}

func (o byIDOption) isStoreOption() {}

func (o byIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) StoreOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isStoreOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."organization_id"`: o.id})
}

func BySlug(slug string) StoreOption {
	return bySlugOption{slug: slug}
}

type bySlugOption struct {
	slug string
}

func (o bySlugOption) isStoreOption() {}

func (o bySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slug})
}

func BySlugs(slugs []string) StoreOption {
	return bySlugsOption{slugs: slugs}
}

type bySlugsOption struct {
	slugs []string
}

func (o bySlugsOption) isStoreOption() {}

func (o bySlugsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slugs})
}

type PageStoreOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isPageStoreOption()
}

func PageByOrganizationID(id uuid.UUID) PageStoreOption {
	return pageByOrganizationIDOption{id: id}
}

type pageByOrganizationIDOption struct {
	id uuid.UUID
}

func (o pageByOrganizationIDOption) isPageStoreOption() {}

func (o pageByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = gp."group_id"`).
		Where(sq.Eq{`g."organization_id"`: o.id})
}

func PageByPageIDs(ids []uuid.UUID) PageStoreOption {
	return pageByPageIDsOption{ids: ids}
}

type pageByPageIDsOption struct {
	ids []uuid.UUID
}

func (o pageByPageIDsOption) isPageStoreOption() {}

func (o pageByPageIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`gp."page_id"`: o.ids})
}

func PageByEnvironmentID(id uuid.UUID) PageStoreOption {
	return pageByEnvironmentIDOption{id: id}
}

type pageByEnvironmentIDOption struct {
	id uuid.UUID
}

func (o pageByEnvironmentIDOption) isPageStoreOption() {}

func (o pageByEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"page" p ON p."id" = gp."page_id"`).
		Where(sq.Eq{`p."environment_id"`: o.id})
}

type Store interface {
	Get(context.Context, ...StoreOption) (*Group, error)
	List(context.Context, ...StoreOption) ([]*Group, error)
	Create(context.Context, *Group) error
	Update(context.Context, *Group) error
	Delete(context.Context, *Group) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)

	ListPages(context.Context, ...PageStoreOption) ([]*GroupPage, error)
	BulkInsertPages(context.Context, []*GroupPage) error
	BulkUpdatePages(context.Context, []*GroupPage) error
	BulkDeletePages(context.Context, []*GroupPage) error
}
