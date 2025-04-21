package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type APIKeyQuery interface{ isAPIKeyQuery() }

type APIKeyByIDQuery struct{ ID uuid.UUID }

func (APIKeyByIDQuery) isAPIKeyQuery() {}

func APIKeyByID(id uuid.UUID) APIKeyQuery { return APIKeyByIDQuery{ID: id} }

type APIKeyByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (APIKeyByOrganizationIDQuery) isAPIKeyQuery() {}

func APIKeyByOrganizationID(organizationID uuid.UUID) APIKeyQuery {
	return APIKeyByOrganizationIDQuery{OrganizationID: organizationID}
}

type APIKeyByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (APIKeyByEnvironmentIDQuery) isAPIKeyQuery() {}

func APIKeyByEnvironmentID(environmentID uuid.UUID) APIKeyQuery {
	return APIKeyByEnvironmentIDQuery{EnvironmentID: environmentID}
}

type APIKeyByEnvironmentIDsQuery struct{ EnvironmentIDs []uuid.UUID }

func (APIKeyByEnvironmentIDsQuery) isAPIKeyQuery() {}

func APIKeyByEnvironmentIDs(environmentIDs []uuid.UUID) APIKeyQuery {
	return APIKeyByEnvironmentIDsQuery{EnvironmentIDs: environmentIDs}
}

type APIKeyByUserIDQuery struct{ UserID uuid.UUID }

func (APIKeyByUserIDQuery) isAPIKeyQuery() {}

func APIKeyByUserID(userID uuid.UUID) APIKeyQuery { return APIKeyByUserIDQuery{UserID: userID} }

type APIKeyByKeyQuery struct{ Key string }

func (APIKeyByKeyQuery) isAPIKeyQuery() {}

func APIKeyByKey(key string) APIKeyQuery { return APIKeyByKeyQuery{Key: key} }

func applyAPIKeyQueries(b sq.SelectBuilder, queries ...APIKeyQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case APIKeyByIDQuery:
			b = b.Where(sq.Eq{`ak."id"`: q.ID})
		case APIKeyByOrganizationIDQuery:
			b = b.Where(sq.Eq{`ak."organization_id"`: q.OrganizationID})
		case APIKeyByEnvironmentIDQuery:
			b = b.Where(sq.Eq{`ak."environment_id"`: q.EnvironmentID})
		case APIKeyByEnvironmentIDsQuery:
			b = b.Where(sq.Eq{`ak."environment_id"`: q.EnvironmentIDs})
		case APIKeyByUserIDQuery:
			b = b.Where(sq.Eq{`ak."user_id"`: q.UserID})
		case APIKeyByKeyQuery:
			b = b.Where(sq.Eq{`ak."key"`: q.Key})
		}
	}
	return b
}
