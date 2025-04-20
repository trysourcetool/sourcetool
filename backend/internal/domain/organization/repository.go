package organization

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Query interface{ isQuery() }

type ByIDQuery struct{ ID uuid.UUID }

func (ByIDQuery) isQuery() {}

func ByID(id uuid.UUID) Query { return ByIDQuery{ID: id} }

type BySubdomainQuery struct{ Subdomain string }

func (BySubdomainQuery) isQuery() {}

func BySubdomain(subdomain string) Query { return BySubdomainQuery{Subdomain: subdomain} }

type ByUserIDQuery struct{ ID uuid.UUID }

func (ByUserIDQuery) isQuery() {}

func ByUserID(id uuid.UUID) Query { return ByUserIDQuery{ID: id} }

type Repository interface {
	Get(context.Context, ...Query) (*Organization, error)
	List(context.Context, ...Query) ([]*Organization, error)
	Create(context.Context, *Organization) error
	IsSubdomainExists(context.Context, string) (bool, error)
}
