package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type GroupOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func GroupByID(id uuid.UUID) GroupOption {
	return groupByIDOption{id: id}
}

type groupByIDOption struct {
	id uuid.UUID
}

func (o groupByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."id"`: o.id})
}

func GroupByOrganizationID(id uuid.UUID) GroupOption {
	return groupByOrganizationIDOption{id: id}
}

type groupByOrganizationIDOption struct {
	id uuid.UUID
}

func (o groupByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."organization_id"`: o.id})
}

func GroupBySlug(slug string) GroupOption {
	return groupBySlugOption{slug: slug}
}

type groupBySlugOption struct {
	slug string
}

func (o groupBySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slug})
}

func GroupBySlugs(slugs []string) GroupOption {
	return groupBySlugsOption{slugs: slugs}
}

type groupBySlugsOption struct {
	slugs []string
}

func (o groupBySlugsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slugs})
}

type GroupPageOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func GroupPageByOrganizationID(id uuid.UUID) GroupPageOption {
	return groupPageByOrganizationIDOption{id: id}
}

type groupPageByOrganizationIDOption struct {
	id uuid.UUID
}

func (o groupPageByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."organization_id"`: o.id})
}

func GroupPageByPageIDs(ids []uuid.UUID) GroupPageOption {
	return groupPageByPageIDsOption{ids: ids}
}

type groupPageByPageIDsOption struct {
	ids []uuid.UUID
}

func (o groupPageByPageIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`gp."page_id"`: o.ids})
}
