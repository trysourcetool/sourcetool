package group

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type RepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isRepositoryOption()
}

func ByID(id uuid.UUID) RepositoryOption {
	return byIDOption{id: id}
}

type byIDOption struct {
	id uuid.UUID
}

func (o byIDOption) isRepositoryOption() {}

func (o byIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) RepositoryOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isRepositoryOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."organization_id"`: o.id})
}

func BySlug(slug string) RepositoryOption {
	return bySlugOption{slug: slug}
}

type bySlugOption struct {
	slug string
}

func (o bySlugOption) isRepositoryOption() {}

func (o bySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slug})
}

func BySlugs(slugs []string) RepositoryOption {
	return bySlugsOption{slugs: slugs}
}

type bySlugsOption struct {
	slugs []string
}

func (o bySlugsOption) isRepositoryOption() {}

func (o bySlugsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slugs})
}

type PageRepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isPageRepositoryOption()
}

func PageByOrganizationID(id uuid.UUID) PageRepositoryOption {
	return pageByOrganizationIDOption{id: id}
}

type pageByOrganizationIDOption struct {
	id uuid.UUID
}

func (o pageByOrganizationIDOption) isPageRepositoryOption() {}

func (o pageByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = gp."group_id"`).
		Where(sq.Eq{`g."organization_id"`: o.id})
}

func PageByPageIDs(ids []uuid.UUID) PageRepositoryOption {
	return pageByPageIDsOption{ids: ids}
}

type pageByPageIDsOption struct {
	ids []uuid.UUID
}

func (o pageByPageIDsOption) isPageRepositoryOption() {}

func (o pageByPageIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`gp."page_id"`: o.ids})
}

func PageByEnvironmentID(id uuid.UUID) PageRepositoryOption {
	return pageByEnvironmentIDOption{id: id}
}

type pageByEnvironmentIDOption struct {
	id uuid.UUID
}

func (o pageByEnvironmentIDOption) isPageRepositoryOption() {}

func (o pageByEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"page" p ON p."id" = gp."page_id"`).
		Where(sq.Eq{`p."environment_id"`: o.id})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*Group, error)
	List(context.Context, ...RepositoryOption) ([]*Group, error)
	Create(context.Context, *Group) error
	Update(context.Context, *Group) error
	Delete(context.Context, *Group) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)

	ListPages(context.Context, ...PageRepositoryOption) ([]*GroupPage, error)
	BulkInsertPages(context.Context, []*GroupPage) error
	BulkUpdatePages(context.Context, []*GroupPage) error
	BulkDeletePages(context.Context, []*GroupPage) error
}
