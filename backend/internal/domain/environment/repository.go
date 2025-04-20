package environment

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

type BySlugQuery struct{ Slug string }

func (BySlugQuery) isQuery() {}

func BySlug(slug string) Query { return BySlugQuery{Slug: slug} }

type Repository interface {
	Get(context.Context, ...Query) (*Environment, error)
	List(context.Context, ...Query) ([]*Environment, error)
	Create(context.Context, *Environment) error
	Update(context.Context, *Environment) error
	Delete(context.Context, *Environment) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)
	BulkInsert(context.Context, []*Environment) error
	MapByAPIKeyIDs(context.Context, []uuid.UUID) (map[uuid.UUID]*Environment, error)
}
