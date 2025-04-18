package page

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type RepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isRepositoryOption()
}

func ByID(id uuid.UUID) RepositoryOption {
	return byIDOption{id: id}
}

type byIDOption struct {
	id uuid.UUID
}

func (o byIDOption) isRepositoryOption() {}

func (o byIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."id"`: o.id})
}

func ByOrganizationID(id uuid.UUID) RepositoryOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isRepositoryOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."organization_id"`: o.id})
}

func ByAPIKeyID(id uuid.UUID) RepositoryOption {
	return byAPIKeyIDOption{id: id}
}

type byAPIKeyIDOption struct {
	id uuid.UUID
}

func (o byAPIKeyIDOption) isRepositoryOption() {}

func (o byAPIKeyIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."api_key_id"`: o.id})
}

func BySessionID(id uuid.UUID) RepositoryOption {
	return bySessionIDOption{id: id}
}

type bySessionIDOption struct {
	id uuid.UUID
}

func (o bySessionIDOption) isRepositoryOption() {}

func (o bySessionIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = p."api_key_id"`).
		InnerJoin(`"session" s ON s."api_key_id" = ak."id"`).
		Where(sq.Eq{`s."id"`: o.id})
}

func ByEnvironmentID(id uuid.UUID) RepositoryOption {
	return byEnvironmentIDOption{id: id}
}

type byEnvironmentIDOption struct {
	id uuid.UUID
}

func (o byEnvironmentIDOption) isRepositoryOption() {}

func (o byEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`p."environment_id"`: o.id})
}

func Limit(limit uint64) RepositoryOption {
	return limitOption{limit: limit}
}

type limitOption struct {
	limit uint64
}

func (o limitOption) isRepositoryOption() {}

func (o limitOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(o.limit)
}

func Offset(offset uint64) RepositoryOption {
	return offsetOption{offset: offset}
}

type offsetOption struct {
	offset uint64
}

func (o offsetOption) isRepositoryOption() {}

func (o offsetOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(o.offset)
}

func OrderBy(orderBy string) RepositoryOption {
	return orderByOption{orderBy: orderBy}
}

type orderByOption struct {
	orderBy string
}

func (o orderByOption) isRepositoryOption() {}

func (o orderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*Page, error)
	List(context.Context, ...RepositoryOption) ([]*Page, error)
	BulkInsert(context.Context, []*Page) error
	BulkUpdate(context.Context, []*Page) error
	BulkDelete(context.Context, []*Page) error
}
