package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type APIKeyOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isAPIKeyOption()
}

func APIKeyByID(id uuid.UUID) APIKeyOption {
	return apiKeyByIDOption{id: id}
}

type apiKeyByIDOption struct {
	id uuid.UUID
}

func (o apiKeyByIDOption) isAPIKeyOption() {}

func (o apiKeyByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."id"`: o.id})
}

func APIKeyByOrganizationID(id uuid.UUID) APIKeyOption {
	return apiKeyByOrganizationIDOption{id: id}
}

type apiKeyByOrganizationIDOption struct {
	id uuid.UUID
}

func (o apiKeyByOrganizationIDOption) isAPIKeyOption() {}

func (o apiKeyByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."organization_id"`: o.id})
}

func APIKeyByEnvironmentID(id uuid.UUID) APIKeyOption {
	return apiKeyByEnvironmentIDOption{id: id}
}

type apiKeyByEnvironmentIDOption struct {
	id uuid.UUID
}

func (o apiKeyByEnvironmentIDOption) isAPIKeyOption() {}

func (o apiKeyByEnvironmentIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: o.id})
}

func APIKeyByEnvironmentIDs(ids []uuid.UUID) APIKeyOption {
	return apiKeyByEnvironmentIDsOption{ids: ids}
}

type apiKeyByEnvironmentIDsOption struct {
	ids []uuid.UUID
}

func (o apiKeyByEnvironmentIDsOption) isAPIKeyOption() {}

func (o apiKeyByEnvironmentIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."environment_id"`: o.ids})
}

func APIKeyByUserID(id uuid.UUID) APIKeyOption {
	return apiKeyByUserIDOption{id: id}
}

type apiKeyByUserIDOption struct {
	id uuid.UUID
}

func (o apiKeyByUserIDOption) isAPIKeyOption() {}

func (o apiKeyByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."user_id"`: o.id})
}

func APIKeyByKey(key string) APIKeyOption {
	return apiKeyByKeyOption{key: key}
}

type apiKeyByKeyOption struct {
	key string
}

func (o apiKeyByKeyOption) isAPIKeyOption() {}

func (o apiKeyByKeyOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ak."key"`: o.key})
}
