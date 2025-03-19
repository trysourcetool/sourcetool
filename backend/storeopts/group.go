package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type GroupOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isGroupOption()
}

func GroupByID(id uuid.UUID) GroupOption {
	return groupByIDOption{id: id}
}

type groupByIDOption struct {
	id uuid.UUID
}

func (o groupByIDOption) isGroupOption() {}

func (o groupByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."id"`: o.id})
}

func GroupByOrganizationID(id uuid.UUID) GroupOption {
	return groupByOrganizationIDOption{id: id}
}

type groupByOrganizationIDOption struct {
	id uuid.UUID
}

func (o groupByOrganizationIDOption) isGroupOption() {}

func (o groupByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."organization_id"`: o.id})
}

func GroupBySlug(slug string) GroupOption {
	return groupBySlugOption{slug: slug}
}

type groupBySlugOption struct {
	slug string
}

func (o groupBySlugOption) isGroupOption() {}

func (o groupBySlugOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slug})
}

func GroupBySlugs(slugs []string) GroupOption {
	return groupBySlugsOption{slugs: slugs}
}

type groupBySlugsOption struct {
	slugs []string
}

func (o groupBySlugsOption) isGroupOption() {}

func (o groupBySlugsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`g."slug"`: o.slugs})
}

type GroupPageOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isGroupPageOption()
}

func GroupPageByOrganizationID(id uuid.UUID) GroupPageOption {
	return groupPageByOrganizationIDOption{id: id}
}

type groupPageByOrganizationIDOption struct {
	id uuid.UUID
}

func (o groupPageByOrganizationIDOption) isGroupPageOption() {}

func (o groupPageByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = gp."group_id"`).
		Where(sq.Eq{`g."organization_id"`: o.id})
}

func GroupPageByPageIDs(ids []uuid.UUID) GroupPageOption {
	return groupPageByPageIDsOption{ids: ids}
}

type groupPageByPageIDsOption struct {
	ids []uuid.UUID
}

func (o groupPageByPageIDsOption) isGroupPageOption() {}

func (o groupPageByPageIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`gp."page_id"`: o.ids})
}

func GroupPageByEnvironmentID(id uuid.UUID) GroupPageOption {
	return groupPageByEnvironmentIDOption{id: id}
}

type groupPageByEnvironmentIDOption struct {
	id uuid.UUID
}

func (o groupPageByEnvironmentIDOption) isGroupPageOption() {}

func (o groupPageByEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"page" p ON p."id" = gp."page_id"`).
		Where(sq.Eq{`p."environment_id"`: o.id})
}
