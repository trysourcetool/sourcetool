//go:build ee
// +build ee

package postgres

import sq "github.com/Masterminds/squirrel"

func applyGroupQueries(b sq.SelectBuilder, queries ...GroupQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case GroupByIDQuery:
			b = b.Where(sq.Eq{`g."id"`: q.ID})
		case GroupByOrganizationIDQuery:
			b = b.Where(sq.Eq{`g."organization_id"`: q.OrganizationID})
		case GroupBySlugQuery:
			b = b.Where(sq.Eq{`g."slug"`: q.Slug})
		case GroupBySlugsQuery:
			b = b.Where(sq.Eq{`g."slug"`: q.Slugs})
		}
	}

	return b
}

func applyGroupPageQueries(b sq.SelectBuilder, queries ...GroupPageQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case GroupPageByOrganizationIDQuery:
			b = b.
				InnerJoin(`"group" g ON g."id" = gp."group_id"`).
				Where(sq.Eq{`g."organization_id"`: q.OrganizationID})
		case GroupPageByPageIDsQuery:
			b = b.Where(sq.Eq{`gp."page_id"`: q.PageIDs})
		case GroupPageByEnvironmentIDQuery:
			b = b.
				InnerJoin(`"page" p ON p."id" = gp."page_id"`).
				Where(sq.Eq{`p."environment_id"`: q.EnvironmentID})
		}
	}

	return b
}
