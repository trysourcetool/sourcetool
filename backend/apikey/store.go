package apikey

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
	return b.Where(sq.Eq{`ak."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) StoreOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isStoreOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."organization_id"`: o.id})
}

func ByEnvironmentID(id uuid.UUID) StoreOption {
	return byEnvironmentIDOption{id: id}
}

type byEnvironmentIDOption struct {
	id uuid.UUID
}

func (o byEnvironmentIDOption) isStoreOption() {}

func (o byEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: o.id})
}

func ByEnvironmentIDs(ids []uuid.UUID) StoreOption {
	return byEnvironmentIDsOption{ids: ids}
}

type byEnvironmentIDsOption struct {
	ids []uuid.UUID
}

func (o byEnvironmentIDsOption) isStoreOption() {}

func (o byEnvironmentIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: o.ids})
}

func ByUserID(id uuid.UUID) StoreOption {
	return byUserIDOption{id: id}
}

type byUserIDOption struct {
	id uuid.UUID
}

func (o byUserIDOption) isStoreOption() {}

func (o byUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."user_id"`: o.id})
}

func ByKey(key string) StoreOption {
	return byKeyOption{key: key}
}

type byKeyOption struct {
	key string
}

func (o byKeyOption) isStoreOption() {}

func (o byKeyOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."key"`: o.key})
}

type Store interface {
	Get(context.Context, ...StoreOption) (*APIKey, error)
	List(context.Context, ...StoreOption) ([]*APIKey, error)
	Create(context.Context, *APIKey) error
	Update(context.Context, *APIKey) error
	Delete(context.Context, *APIKey) error
}
