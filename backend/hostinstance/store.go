package hostinstance

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
	return b.Where(sq.Eq{`hi."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) StoreOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isStoreOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."organization_id"`: o.id})
}

func ByAPIKeyID(id uuid.UUID) StoreOption {
	return byAPIKeyIDOption{id: id}
}

type byAPIKeyIDOption struct {
	id uuid.UUID
}

func (o byAPIKeyIDOption) isStoreOption() {}

func (o byAPIKeyIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."api_key_id"`: o.id})
}

func ByAPIKey(key string) StoreOption {
	return byAPIKeyOption{key: key}
}

type byAPIKeyOption struct {
	key string
}

func (o byAPIKeyOption) isStoreOption() {}

func (o byAPIKeyOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = hi."api_key_id"`).
		Where(sq.Eq{`ak."key"`: o.key})
}

type Store interface {
	Get(context.Context, ...StoreOption) (*HostInstance, error)
	List(context.Context, ...StoreOption) ([]*HostInstance, error)
	Create(context.Context, *HostInstance) error
	Update(context.Context, *HostInstance) error
}
