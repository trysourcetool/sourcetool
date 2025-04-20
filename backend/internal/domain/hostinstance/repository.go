package hostinstance

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Query interface{ isQuery() }

type ByIDQuery struct{ ID uuid.UUID }

func (ByIDQuery) isQuery() {}

func ByID(id uuid.UUID) Query { return ByIDQuery{ID: id} }

type ByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (ByOrganizationIDQuery) isQuery() {}

func ByOrganizationID(organizationID uuid.UUID) Query {
	return ByOrganizationIDQuery{OrganizationID: organizationID}
}

type ByAPIKeyIDQuery struct{ APIKeyID uuid.UUID }

func (ByAPIKeyIDQuery) isQuery() {}

func ByAPIKeyID(apiKeyID uuid.UUID) Query { return ByAPIKeyIDQuery{APIKeyID: apiKeyID} }

type ByAPIKeyQuery struct{ APIKey string }

func (ByAPIKeyQuery) isQuery() {}

func ByAPIKey(apiKey string) Query { return ByAPIKeyQuery{APIKey: apiKey} }

type Repository interface {
	Get(context.Context, ...Query) (*HostInstance, error)
	List(context.Context, ...Query) ([]*HostInstance, error)
	Create(context.Context, *HostInstance) error
	Update(context.Context, *HostInstance) error
}
