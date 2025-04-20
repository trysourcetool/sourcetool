package page

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Query interface{ isQuery() }

type ByIDQuery struct{ ID uuid.UUID }

func (q ByIDQuery) isQuery() {}

func ByID(id uuid.UUID) Query { return ByIDQuery{ID: id} }

type ByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q ByOrganizationIDQuery) isQuery() {}

func ByOrganizationID(id uuid.UUID) Query { return ByOrganizationIDQuery{OrganizationID: id} }

type ByAPIKeyIDQuery struct{ APIKeyID uuid.UUID }

func (q ByAPIKeyIDQuery) isQuery() {}

func ByAPIKeyID(id uuid.UUID) Query { return ByAPIKeyIDQuery{APIKeyID: id} }

type BySessionIDQuery struct{ SessionID uuid.UUID }

func (q BySessionIDQuery) isQuery() {}

func BySessionID(id uuid.UUID) Query { return BySessionIDQuery{SessionID: id} }

type ByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (q ByEnvironmentIDQuery) isQuery() {}

func ByEnvironmentID(id uuid.UUID) Query { return ByEnvironmentIDQuery{EnvironmentID: id} }

type LimitQuery struct{ Limit uint64 }

func (q LimitQuery) isQuery() {}

func Limit(limit uint64) Query { return LimitQuery{Limit: limit} }

type OffsetQuery struct{ Offset uint64 }

func (q OffsetQuery) isQuery() {}

func Offset(offset uint64) Query { return OffsetQuery{Offset: offset} }

type OrderByQuery struct{ OrderBy string }

func (q OrderByQuery) isQuery() {}

func OrderBy(orderBy string) Query { return OrderByQuery{OrderBy: orderBy} }

type Repository interface {
	Get(context.Context, ...Query) (*Page, error)
	List(context.Context, ...Query) ([]*Page, error)
	BulkInsert(context.Context, []*Page) error
	BulkUpdate(context.Context, []*Page) error
	BulkDelete(context.Context, []*Page) error
}
