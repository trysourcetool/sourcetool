package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type HostInstanceQuery interface{ isHostInstanceQuery() }

type HostInstanceByIDQuery struct{ ID uuid.UUID }

func (HostInstanceByIDQuery) isHostInstanceQuery() {}

func HostInstanceByID(id uuid.UUID) HostInstanceQuery { return HostInstanceByIDQuery{ID: id} }

type HostInstanceByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (HostInstanceByOrganizationIDQuery) isHostInstanceQuery() {}

func HostInstanceByOrganizationID(organizationID uuid.UUID) HostInstanceQuery {
	return HostInstanceByOrganizationIDQuery{OrganizationID: organizationID}
}

type HostInstanceByAPIKeyIDQuery struct{ APIKeyID uuid.UUID }

func (HostInstanceByAPIKeyIDQuery) isHostInstanceQuery() {}

func HostInstanceByAPIKeyID(apiKeyID uuid.UUID) HostInstanceQuery {
	return HostInstanceByAPIKeyIDQuery{APIKeyID: apiKeyID}
}

type HostInstanceByAPIKeyQuery struct{ APIKey string }

func (HostInstanceByAPIKeyQuery) isHostInstanceQuery() {}

func HostInstanceByAPIKey(apiKey string) HostInstanceQuery {
	return HostInstanceByAPIKeyQuery{APIKey: apiKey}
}

type HostInstanceBySessionIDQuery struct{ SessionID uuid.UUID }

func (HostInstanceBySessionIDQuery) isHostInstanceQuery() {}

func HostInstanceBySessionID(sessionID uuid.UUID) HostInstanceQuery {
	return HostInstanceBySessionIDQuery{SessionID: sessionID}
}

type HostInstanceStore interface {
	Get(ctx context.Context, queries ...HostInstanceQuery) (*core.HostInstance, error)
	List(ctx context.Context, queries ...HostInstanceQuery) ([]*core.HostInstance, error)
	Create(ctx context.Context, m *core.HostInstance) error
	Update(ctx context.Context, m *core.HostInstance) error
}
