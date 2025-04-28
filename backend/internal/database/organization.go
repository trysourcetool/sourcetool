package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type OrganizationQuery interface{ isOrganizationQuery() }

type OrganizationByIDQuery struct{ ID uuid.UUID }

func (OrganizationByIDQuery) isOrganizationQuery() {}

func OrganizationByID(id uuid.UUID) OrganizationQuery { return OrganizationByIDQuery{ID: id} }

type OrganizationBySubdomainQuery struct{ Subdomain string }

func (OrganizationBySubdomainQuery) isOrganizationQuery() {}

func OrganizationBySubdomain(subdomain string) OrganizationQuery {
	return OrganizationBySubdomainQuery{Subdomain: subdomain}
}

type OrganizationByUserIDQuery struct{ ID uuid.UUID }

func (OrganizationByUserIDQuery) isOrganizationQuery() {}

func OrganizationByUserID(id uuid.UUID) OrganizationQuery {
	return OrganizationByUserIDQuery{ID: id}
}

type OrganizationStore interface {
	Get(ctx context.Context, queries ...OrganizationQuery) (*core.Organization, error)
	List(ctx context.Context, queries ...OrganizationQuery) ([]*core.Organization, error)
	Create(ctx context.Context, m *core.Organization) error
	IsSubdomainExists(ctx context.Context, subdomain string) (bool, error)
}
