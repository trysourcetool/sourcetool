package hostinstance

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
	return b.Where(sq.Eq{`hi."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) RepositoryOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isRepositoryOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."organization_id"`: o.id})
}

func ByAPIKeyID(id uuid.UUID) RepositoryOption {
	return byAPIKeyIDOption{id: id}
}

type byAPIKeyIDOption struct {
	id uuid.UUID
}

func (o byAPIKeyIDOption) isRepositoryOption() {}

func (o byAPIKeyIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."api_key_id"`: o.id})
}

func ByAPIKey(key string) RepositoryOption {
	return byAPIKeyOption{key: key}
}

type byAPIKeyOption struct {
	key string
}

func (o byAPIKeyOption) isRepositoryOption() {}

func (o byAPIKeyOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = hi."api_key_id"`).
		Where(sq.Eq{`ak."key"`: o.key})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*HostInstance, error)
	List(context.Context, ...RepositoryOption) ([]*HostInstance, error)
	Create(context.Context, *HostInstance) error
	Update(context.Context, *HostInstance) error
}
