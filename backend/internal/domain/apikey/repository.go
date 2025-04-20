package apikey

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

type ByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (ByEnvironmentIDQuery) isQuery() {}

func ByEnvironmentID(environmentID uuid.UUID) Query {
	return ByEnvironmentIDQuery{EnvironmentID: environmentID}
}

type ByEnvironmentIDsQuery struct{ EnvironmentIDs []uuid.UUID }

func (ByEnvironmentIDsQuery) isQuery() {}

func ByEnvironmentIDs(environmentIDs []uuid.UUID) Query {
	return ByEnvironmentIDsQuery{EnvironmentIDs: environmentIDs}
}

type ByUserIDQuery struct{ UserID uuid.UUID }

func (ByUserIDQuery) isQuery() {}

func ByUserID(userID uuid.UUID) Query { return ByUserIDQuery{UserID: userID} }

type ByKeyQuery struct{ Key string }

func (ByKeyQuery) isQuery() {}

func ByKey(key string) Query { return ByKeyQuery{Key: key} }

type Repository interface {
	Get(context.Context, ...Query) (*APIKey, error)
	List(context.Context, ...Query) ([]*APIKey, error)
	Create(context.Context, *APIKey) error
	Update(context.Context, *APIKey) error
	Delete(context.Context, *APIKey) error
}
