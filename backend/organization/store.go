package organization

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
	return b.Where(sq.Eq{`o."id"`: o.id})
}

func BySubdomain(subdomain string) StoreOption {
	return bySubdomainOption{subdomain: subdomain}
}

type bySubdomainOption struct {
	subdomain string
}

func (o bySubdomainOption) isStoreOption() {}

func (o bySubdomainOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`o."subdomain"`: o.subdomain})
}

func ByUserID(id uuid.UUID) StoreOption {
	return byUserIDOption{id: id}
}

type byUserIDOption struct {
	id uuid.UUID
}

func (o byUserIDOption) isStoreOption() {}

func (o byUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
		Where(sq.Eq{`uoa."user_id"`: o.id})
}

type Store interface {
	Get(context.Context, ...StoreOption) (*Organization, error)
	List(context.Context, ...StoreOption) ([]*Organization, error)
	Create(context.Context, *Organization) error
	IsSubdomainExists(context.Context, string) (bool, error)
}
