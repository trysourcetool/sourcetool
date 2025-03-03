package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type EnvironmentOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func EnvironmentByID(id uuid.UUID) EnvironmentOption {
	return environmentByIDOption{id: id}
}

type environmentByIDOption struct {
	id uuid.UUID
}

func (o environmentByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."id"`: o.id})
}

func EnvironmentByOrganizationID(id uuid.UUID) EnvironmentOption {
	return environmentByOrganizationIDOption{id: id}
}

type environmentByOrganizationIDOption struct {
	id uuid.UUID
}

func (o environmentByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."organization_id"`: o.id})
}

func EnvironmentBySlug(slug string) EnvironmentOption {
	return environmentBySlugOption{slug: slug}
}

type environmentBySlugOption struct {
	slug string
}

func (o environmentBySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`e."slug"`: o.slug})
}
