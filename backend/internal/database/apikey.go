package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type APIKeyQuery interface{ isAPIKeyQuery() }

type APIKeyByIDQuery struct{ ID uuid.UUID }

func (APIKeyByIDQuery) isAPIKeyQuery() {}

func APIKeyByID(id uuid.UUID) APIKeyQuery { return APIKeyByIDQuery{ID: id} }

type APIKeyByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (APIKeyByOrganizationIDQuery) isAPIKeyQuery() {}

func APIKeyByOrganizationID(organizationID uuid.UUID) APIKeyQuery {
	return APIKeyByOrganizationIDQuery{OrganizationID: organizationID}
}

type APIKeyByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (APIKeyByEnvironmentIDQuery) isAPIKeyQuery() {}

func APIKeyByEnvironmentID(environmentID uuid.UUID) APIKeyQuery {
	return APIKeyByEnvironmentIDQuery{EnvironmentID: environmentID}
}

type APIKeyByEnvironmentIDsQuery struct{ EnvironmentIDs []uuid.UUID }

func (APIKeyByEnvironmentIDsQuery) isAPIKeyQuery() {}

func APIKeyByEnvironmentIDs(environmentIDs []uuid.UUID) APIKeyQuery {
	return APIKeyByEnvironmentIDsQuery{EnvironmentIDs: environmentIDs}
}

type APIKeyByUserIDQuery struct{ UserID uuid.UUID }

func (APIKeyByUserIDQuery) isAPIKeyQuery() {}

func APIKeyByUserID(userID uuid.UUID) APIKeyQuery {
	return APIKeyByUserIDQuery{UserID: userID}
}

type APIKeyByKeyHashQuery struct{ KeyHash string }

func (APIKeyByKeyHashQuery) isAPIKeyQuery() {}

func APIKeyByKeyHash(keyHash string) APIKeyQuery { return APIKeyByKeyHashQuery{KeyHash: keyHash} }

type APIKeyStore interface {
	Get(ctx context.Context, queries ...APIKeyQuery) (*core.APIKey, error)
	List(ctx context.Context, queries ...APIKeyQuery) ([]*core.APIKey, error)
	Create(ctx context.Context, m *core.APIKey) error
	Update(ctx context.Context, m *core.APIKey) error
	Delete(ctx context.Context, m *core.APIKey) error
}
