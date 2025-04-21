package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type PageQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isPageQuery()
}

type pageByIDQuery struct{ id uuid.UUID }

func (q pageByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."id"`: q.id})
}

func (pageByIDQuery) isPageQuery() {}

func PageByID(id uuid.UUID) PageQuery { return pageByIDQuery{id: id} }

type pageByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q pageByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."organization_id"`: q.organizationID})
}

func (pageByOrganizationIDQuery) isPageQuery() {}

func PageByOrganizationID(id uuid.UUID) PageQuery {
	return pageByOrganizationIDQuery{organizationID: id}
}

type pageByAPIKeyIDQuery struct{ apiKeyID uuid.UUID }

func (q pageByAPIKeyIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."api_key_id"`: q.apiKeyID})
}

func (pageByAPIKeyIDQuery) isPageQuery() {}

func PageByAPIKeyID(id uuid.UUID) PageQuery { return pageByAPIKeyIDQuery{apiKeyID: id} }

type pageBySessionIDQuery struct{ sessionID uuid.UUID }

func (q pageBySessionIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = p."api_key_id"`).
		InnerJoin(`"session" s ON s."api_key_id" = ak."id"`).
		Where(sq.Eq{`s."id"`: q.sessionID})
}

func (pageBySessionIDQuery) isPageQuery() {}

func PageBySessionID(id uuid.UUID) PageQuery { return pageBySessionIDQuery{sessionID: id} }

type pageByEnvironmentIDQuery struct{ environmentID uuid.UUID }

func (q pageByEnvironmentIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."environment_id"`: q.environmentID})
}

func (pageByEnvironmentIDQuery) isPageQuery() {}

func PageByEnvironmentID(id uuid.UUID) PageQuery { return pageByEnvironmentIDQuery{environmentID: id} }

type pageLimitQuery struct{ limit uint64 }

func (q pageLimitQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(q.limit)
}

func (pageLimitQuery) isPageQuery() {}

func PageLimit(limit uint64) PageQuery { return pageLimitQuery{limit: limit} }

type pageOffsetQuery struct{ offset uint64 }

func (q pageOffsetQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(q.offset)
}

func (pageOffsetQuery) isPageQuery() {}

func PageOffset(offset uint64) PageQuery { return pageOffsetQuery{offset: offset} }

type pageOrderByQuery struct{ orderBy string }

func (q pageOrderByQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(q.orderBy)
}

func (pageOrderByQuery) isPageQuery() {}

func PageOrderBy(orderBy string) PageQuery { return pageOrderByQuery{orderBy: orderBy} }
