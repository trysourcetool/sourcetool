package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type APIKeyQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isAPIKeyQuery()
}

type apiKeyByIDQuery struct{ id uuid.UUID }

func (q apiKeyByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."id"`: q.id})
}

func (apiKeyByIDQuery) isAPIKeyQuery() {}

func APIKeyByID(id uuid.UUID) APIKeyQuery { return apiKeyByIDQuery{id: id} }

type apiKeyByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q apiKeyByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."organization_id"`: q.organizationID})
}

func (apiKeyByOrganizationIDQuery) isAPIKeyQuery() {}

func APIKeyByOrganizationID(organizationID uuid.UUID) APIKeyQuery {
	return apiKeyByOrganizationIDQuery{organizationID: organizationID}
}

type apiKeyByEnvironmentIDQuery struct{ environmentID uuid.UUID }

func (q apiKeyByEnvironmentIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: q.environmentID})
}

func (apiKeyByEnvironmentIDQuery) isAPIKeyQuery() {}

func APIKeyByEnvironmentID(environmentID uuid.UUID) APIKeyQuery {
	return apiKeyByEnvironmentIDQuery{environmentID: environmentID}
}

type apiKeyByEnvironmentIDsQuery struct{ environmentIDs []uuid.UUID }

func (q apiKeyByEnvironmentIDsQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: q.environmentIDs})
}

func (apiKeyByEnvironmentIDsQuery) isAPIKeyQuery() {}

func APIKeyByEnvironmentIDs(environmentIDs []uuid.UUID) APIKeyQuery {
	return apiKeyByEnvironmentIDsQuery{environmentIDs: environmentIDs}
}

type apiKeyByUserIDQuery struct{ userID uuid.UUID }

func (q apiKeyByUserIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."user_id"`: q.userID})
}

func (q apiKeyByUserIDQuery) isAPIKeyQuery() {}

func APIKeyByUserID(userID uuid.UUID) APIKeyQuery { return apiKeyByUserIDQuery{userID: userID} }

type apiKeyByKeyQuery struct{ key string }

func (q apiKeyByKeyQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."key"`: q.key})
}

func (apiKeyByKeyQuery) isAPIKeyQuery() {}

func APIKeyByKey(key string) APIKeyQuery { return apiKeyByKeyQuery{key: key} }
