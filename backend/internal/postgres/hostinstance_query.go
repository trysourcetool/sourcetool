package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type HostInstanceQuery interface{ isHostInstanceQuery() }

type HostInstanceByIDQuery struct{ ID uuid.UUID }

func (HostInstanceByIDQuery) isHostInstanceQuery() {}

func HostInstanceByID(id uuid.UUID) HostInstanceQuery { return HostInstanceByIDQuery{ID: id} }

type HostInstanceByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (HostInstanceByOrganizationIDQuery) isHostInstanceQuery() {}

func HostInstanceByOrganizationID(organizationID uuid.UUID) HostInstanceQuery {
	return HostInstanceByOrganizationIDQuery{OrganizationID: organizationID}
}

type HostInstanceByAPIKeyIDQuery struct{ APIKeyID uuid.UUID }

func (HostInstanceByAPIKeyIDQuery) isHostInstanceQuery() {}

func HostInstanceByAPIKeyID(apiKeyID uuid.UUID) HostInstanceQuery {
	return HostInstanceByAPIKeyIDQuery{APIKeyID: apiKeyID}
}

type HostInstanceByAPIKeyQuery struct{ APIKey string }

func (HostInstanceByAPIKeyQuery) isHostInstanceQuery() {}

func HostInstanceByAPIKey(apiKey string) HostInstanceQuery {
	return HostInstanceByAPIKeyQuery{APIKey: apiKey}
}

func applyHostInstanceQueries(b sq.SelectBuilder, queries ...HostInstanceQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case HostInstanceByIDQuery:
			b = b.Where(sq.Eq{`hi."id"`: q.ID})
		case HostInstanceByOrganizationIDQuery:
			b = b.Where(sq.Eq{`hi."organization_id"`: q.OrganizationID})
		case HostInstanceByAPIKeyIDQuery:
			b = b.Where(sq.Eq{`hi."api_key_id"`: q.APIKeyID})
		case HostInstanceByAPIKeyQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."id" = hi."api_key_id"`).
				Where(sq.Eq{`ak."key"`: q.APIKey})
		}
	}

	return b
}
