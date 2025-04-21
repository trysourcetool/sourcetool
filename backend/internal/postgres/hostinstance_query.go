package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type HostInstanceQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isHostInstanceQuery()
}

type hostInstanceByIDQuery struct{ id uuid.UUID }

func (q hostInstanceByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."id"`: q.id})
}

func (hostInstanceByIDQuery) isHostInstanceQuery() {}

func HostInstanceByID(id uuid.UUID) HostInstanceQuery { return hostInstanceByIDQuery{id: id} }

type hostInstanceByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q hostInstanceByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."organization_id"`: q.organizationID})
}

func (hostInstanceByOrganizationIDQuery) isHostInstanceQuery() {}

func HostInstanceByOrganizationID(organizationID uuid.UUID) HostInstanceQuery {
	return hostInstanceByOrganizationIDQuery{organizationID: organizationID}
}

type hostInstanceByAPIKeyIDQuery struct{ apiKeyID uuid.UUID }

func (q hostInstanceByAPIKeyIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."api_key_id"`: q.apiKeyID})
}

func (hostInstanceByAPIKeyIDQuery) isHostInstanceQuery() {}

func HostInstanceByAPIKeyID(apiKeyID uuid.UUID) HostInstanceQuery {
	return hostInstanceByAPIKeyIDQuery{apiKeyID: apiKeyID}
}

type hostInstanceByKeyQuery struct{ key string }

func (q hostInstanceByKeyQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = hi."api_key_id"`).
		Where(sq.Eq{`ak."key"`: q.key})
}

func (hostInstanceByKeyQuery) isHostInstanceQuery() {}

func HostInstanceByKey(key string) HostInstanceQuery {
	return hostInstanceByKeyQuery{key: key}
}
