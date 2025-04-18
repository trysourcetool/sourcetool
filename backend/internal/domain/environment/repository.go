package environment

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
	return b.Where(sq.Eq{`e."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) RepositoryOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isRepositoryOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."organization_id"`: o.id})
}

func BySlug(slug string) RepositoryOption {
	return bySlugOption{slug: slug}
}

type bySlugOption struct {
	slug string
}

func (o bySlugOption) isRepositoryOption() {}

func (o bySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."slug"`: o.slug})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*Environment, error)
	List(context.Context, ...RepositoryOption) ([]*Environment, error)
	Create(context.Context, *Environment) error
	Update(context.Context, *Environment) error
	Delete(context.Context, *Environment) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)
	BulkInsert(context.Context, []*Environment) error
	MapByAPIKeyIDs(context.Context, []uuid.UUID) (map[uuid.UUID]*Environment, error)
}
