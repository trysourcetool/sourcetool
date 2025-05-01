package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type EnvironmentQuery interface{ isEnvironmentQuery() }

type EnvironmentByIDQuery struct{ ID uuid.UUID }

func (EnvironmentByIDQuery) isEnvironmentQuery() {}

func EnvironmentByID(id uuid.UUID) EnvironmentQuery { return EnvironmentByIDQuery{ID: id} }

type EnvironmentByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (EnvironmentByOrganizationIDQuery) isEnvironmentQuery() {}

func EnvironmentByOrganizationID(organizationID uuid.UUID) EnvironmentQuery {
	return EnvironmentByOrganizationIDQuery{OrganizationID: organizationID}
}

type EnvironmentBySlugQuery struct{ Slug string }

func (EnvironmentBySlugQuery) isEnvironmentQuery() {}

func EnvironmentBySlug(slug string) EnvironmentQuery { return EnvironmentBySlugQuery{Slug: slug} }

type EnvironmentByAPIKeyIDsQuery struct{ APIKeyIDs []uuid.UUID }

func (EnvironmentByAPIKeyIDsQuery) isEnvironmentQuery() {}

func EnvironmentByAPIKeyIDs(ids []uuid.UUID) EnvironmentQuery {
	return EnvironmentByAPIKeyIDsQuery{APIKeyIDs: ids}
}

type EnvironmentStore interface {
	Get(ctx context.Context, queries ...EnvironmentQuery) (*core.Environment, error)
	List(ctx context.Context, queries ...EnvironmentQuery) ([]*core.Environment, error)
	Create(ctx context.Context, m *core.Environment) error
	Update(ctx context.Context, m *core.Environment) error
	Delete(ctx context.Context, m *core.Environment) error
	BulkInsert(ctx context.Context, m []*core.Environment) error
	MapByAPIKeyIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*core.Environment, error)
	IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error)
}
