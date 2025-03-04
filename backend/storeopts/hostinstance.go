package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type HostInstanceOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isHostInstanceOption()
}

func HostInstanceByID(id uuid.UUID) HostInstanceOption {
	return hostInstanceByIDOption{id: id}
}

type hostInstanceByIDOption struct {
	id uuid.UUID
}

func (o hostInstanceByIDOption) isHostInstanceOption() {}

func (o hostInstanceByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."id"`: o.id})
}

func HostInstanceByOrganizationID(id uuid.UUID) HostInstanceOption {
	return hostInstanceByOrganizationIDOption{id: id}
}

type hostInstanceByOrganizationIDOption struct {
	id uuid.UUID
}

func (o hostInstanceByOrganizationIDOption) isHostInstanceOption() {}

func (o hostInstanceByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."organization_id"`: o.id})
}

func HostInstanceByAPIKeyID(id uuid.UUID) HostInstanceOption {
	return hostInstanceByAPIKeyIDOption{id: id}
}

type hostInstanceByAPIKeyIDOption struct {
	id uuid.UUID
}

func (o hostInstanceByAPIKeyIDOption) isHostInstanceOption() {}

func (o hostInstanceByAPIKeyIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`hi."api_key_id"`: o.id})
}

func HostInstanceByAPIKey(key string) HostInstanceOption {
	return hostInstanceByAPIKeyOption{key: key}
}

type hostInstanceByAPIKeyOption struct {
	key string
}

func (o hostInstanceByAPIKeyOption) isHostInstanceOption() {}

func (o hostInstanceByAPIKeyOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"api_key" ak ON ak."id" = hi."api_key_id"`).
		Where(sq.Eq{`ak."key"`: o.key})
}
