package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type PageQuery interface{ isPageQuery() }

type PageByIDQuery struct{ ID uuid.UUID }

func (q PageByIDQuery) isPageQuery() {}

func PageByID(id uuid.UUID) PageQuery { return PageByIDQuery{ID: id} }

type PageByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q PageByOrganizationIDQuery) isPageQuery() {}

func PageByOrganizationID(id uuid.UUID) PageQuery {
	return PageByOrganizationIDQuery{OrganizationID: id}
}

type PageByAPIKeyIDQuery struct{ APIKeyID uuid.UUID }

func (q PageByAPIKeyIDQuery) isPageQuery() {}

func PageByAPIKeyID(id uuid.UUID) PageQuery { return PageByAPIKeyIDQuery{APIKeyID: id} }

type PageBySessionIDQuery struct{ SessionID uuid.UUID }

func (q PageBySessionIDQuery) isPageQuery() {}

func PageBySessionID(id uuid.UUID) PageQuery { return PageBySessionIDQuery{SessionID: id} }

type PageByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (q PageByEnvironmentIDQuery) isPageQuery() {}

func PageByEnvironmentID(id uuid.UUID) PageQuery {
	return PageByEnvironmentIDQuery{EnvironmentID: id}
}

type PageLimitQuery struct{ Limit uint64 }

func (q PageLimitQuery) isPageQuery() {}

func PageLimit(limit uint64) PageQuery { return PageLimitQuery{Limit: limit} }

type PageOffsetQuery struct{ Offset uint64 }

func (q PageOffsetQuery) isPageQuery() {}

func PageOffset(offset uint64) PageQuery { return PageOffsetQuery{Offset: offset} }

type PageOrderByQuery struct{ OrderBy string }

func (q PageOrderByQuery) isPageQuery() {}

func PageOrderBy(orderBy string) PageQuery { return PageOrderByQuery{OrderBy: orderBy} }

type PageStore interface {
	Get(ctx context.Context, queries ...PageQuery) (*core.Page, error)
	List(ctx context.Context, queries ...PageQuery) ([]*core.Page, error)
	BulkInsert(ctx context.Context, m []*core.Page) error
	BulkUpdate(ctx context.Context, m []*core.Page) error
	BulkDelete(ctx context.Context, m []*core.Page) error
}
