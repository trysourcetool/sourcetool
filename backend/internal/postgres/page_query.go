package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
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

func PageByEnvironmentID(id uuid.UUID) PageQuery { return PageByEnvironmentIDQuery{EnvironmentID: id} }

type PageLimitQuery struct{ Limit uint64 }

func (q PageLimitQuery) isPageQuery() {}

func PageLimit(limit uint64) PageQuery { return PageLimitQuery{Limit: limit} }

type PageOffsetQuery struct{ Offset uint64 }

func (q PageOffsetQuery) isPageQuery() {}

func PageOffset(offset uint64) PageQuery { return PageOffsetQuery{Offset: offset} }

type PageOrderByQuery struct{ OrderBy string }

func (q PageOrderByQuery) isPageQuery() {}

func PageOrderBy(orderBy string) PageQuery { return PageOrderByQuery{OrderBy: orderBy} }

func applyPageQueries(b sq.SelectBuilder, queries ...PageQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case PageByIDQuery:
			b = b.Where(sq.Eq{`p."id"`: q.ID})
		case PageByOrganizationIDQuery:
			b = b.Where(sq.Eq{`p."organization_id"`: q.OrganizationID})
		case PageByAPIKeyIDQuery:
			b = b.Where(sq.Eq{`p."api_key_id"`: q.APIKeyID})
		case PageBySessionIDQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."id" = p."api_key_id"`).
				InnerJoin(`"session" s ON s."api_key_id" = ak."id"`).
				Where(sq.Eq{`s."id"`: q.SessionID})
		case PageByEnvironmentIDQuery:
			b = b.Where(sq.Eq{`p."environment_id"`: q.EnvironmentID})
		case PageLimitQuery:
			b = b.Limit(q.Limit)
		case PageOffsetQuery:
			b = b.Offset(q.Offset)
		case PageOrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}

	return b
}
