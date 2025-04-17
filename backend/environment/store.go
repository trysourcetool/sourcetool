package environment

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
	return b.Where(sq.Eq{`e."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) StoreOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isStoreOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."organization_id"`: o.id})
}

func BySlug(slug string) StoreOption {
	return bySlugOption{slug: slug}
}

type bySlugOption struct {
	slug string
}

func (o bySlugOption) isStoreOption() {}

func (o bySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."slug"`: o.slug})
}

type Store interface {
	Get(context.Context, ...StoreOption) (*Environment, error)
	List(context.Context, ...StoreOption) ([]*Environment, error)
	Create(context.Context, *Environment) error
	Update(context.Context, *Environment) error
	Delete(context.Context, *Environment) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)
	BulkInsert(context.Context, []*Environment) error
	MapByAPIKeyIDs(context.Context, []uuid.UUID) (map[uuid.UUID]*Environment, error)
}
