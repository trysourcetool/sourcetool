package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type PageOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isPageOption()
}

func PageByID(id uuid.UUID) PageOption {
	return pageByIDOption{id: id}
}

type pageByIDOption struct {
	id uuid.UUID
}

func (o pageByIDOption) isPageOption() {}

func (o pageByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."id"`: o.id})
}

func PageByOrganizationID(id uuid.UUID) PageOption {
	return pageByOrganizationIDOption{id: id}
}

type pageByOrganizationIDOption struct {
	id uuid.UUID
}

func (o pageByOrganizationIDOption) isPageOption() {}

func (o pageByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."organization_id"`: o.id})
}

func PageByAPIKeyID(id uuid.UUID) PageOption {
	return pageByAPIKeyIDOption{id: id}
}

type pageByAPIKeyIDOption struct {
	id uuid.UUID
}

func (o pageByAPIKeyIDOption) isPageOption() {}

func (o pageByAPIKeyIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."api_key_id"`: o.id})
}

func PageBySessionID(id uuid.UUID) PageOption {
	return pageBySessionIDOption{id: id}
}

type pageBySessionIDOption struct {
	id uuid.UUID
}

func (o pageBySessionIDOption) isPageOption() {}

func (o pageBySessionIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = p."api_key_id"`).
		InnerJoin(`"session" s ON s."api_key_id" = ak."id"`).
		Where(sq.Eq{`s."id"`: o.id})
}

func PageByEnvironmentID(id uuid.UUID) PageOption {
	return pageByEnvironmentIDOption{id: id}
}

type pageByEnvironmentIDOption struct {
	id uuid.UUID
}

func (o pageByEnvironmentIDOption) isPageOption() {}

func (o pageByEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."environment_id"`: o.id})
}

func PageLimit(limit uint64) PageOption {
	return pageLimitOption{limit: limit}
}

type pageLimitOption struct {
	limit uint64
}

func (o pageLimitOption) isPageOption() {}

func (o pageLimitOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(o.limit)
}

func PageOffset(offset uint64) PageOption {
	return pageOffsetOption{offset: offset}
}

type pageOffsetOption struct {
	offset uint64
}

func (o pageOffsetOption) isPageOption() {}

func (o pageOffsetOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(o.offset)
}

func PageOrderBy(orderBy string) PageOption {
	return pageOrderByOption{orderBy: orderBy}
}

type pageOrderByOption struct {
	orderBy string
}

func (o pageOrderByOption) isPageOption() {}

func (o pageOrderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}
