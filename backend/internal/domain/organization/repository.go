package organization

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
	return b.Where(sq.Eq{`o."id"`: o.id})
}

func BySubdomain(subdomain string) RepositoryOption {
	return bySubdomainOption{subdomain: subdomain}
}

type bySubdomainOption struct {
	subdomain string
}

func (o bySubdomainOption) isRepositoryOption() {}

func (o bySubdomainOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`o."subdomain"`: o.subdomain})
}

func ByUserID(id uuid.UUID) RepositoryOption {
	return byUserIDOption{id: id}
}

type byUserIDOption struct {
	id uuid.UUID
}

func (o byUserIDOption) isRepositoryOption() {}

func (o byUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
		Where(sq.Eq{`uoa."user_id"`: o.id})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*Organization, error)
	List(context.Context, ...RepositoryOption) ([]*Organization, error)
	Create(context.Context, *Organization) error
	IsSubdomainExists(context.Context, string) (bool, error)
}
