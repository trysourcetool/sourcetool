package apikey

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
	return b.Where(sq.Eq{`ak."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) RepositoryOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isRepositoryOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."organization_id"`: o.id})
}

func ByEnvironmentID(id uuid.UUID) RepositoryOption {
	return byEnvironmentIDOption{id: id}
}

type byEnvironmentIDOption struct {
	id uuid.UUID
}

func (o byEnvironmentIDOption) isRepositoryOption() {}

func (o byEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: o.id})
}

func ByEnvironmentIDs(ids []uuid.UUID) RepositoryOption {
	return byEnvironmentIDsOption{ids: ids}
}

type byEnvironmentIDsOption struct {
	ids []uuid.UUID
}

func (o byEnvironmentIDsOption) isRepositoryOption() {}

func (o byEnvironmentIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: o.ids})
}

func ByUserID(id uuid.UUID) RepositoryOption {
	return byUserIDOption{id: id}
}

type byUserIDOption struct {
	id uuid.UUID
}

func (o byUserIDOption) isRepositoryOption() {}

func (o byUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."user_id"`: o.id})
}

func ByKey(key string) RepositoryOption {
	return byKeyOption{key: key}
}

type byKeyOption struct {
	key string
}

func (o byKeyOption) isRepositoryOption() {}

func (o byKeyOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."key"`: o.key})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*APIKey, error)
	List(context.Context, ...RepositoryOption) ([]*APIKey, error)
	Create(context.Context, *APIKey) error
	Update(context.Context, *APIKey) error
	Delete(context.Context, *APIKey) error
}
