package database

import (
	"context"
)

type Stores interface {
	APIKey() APIKeyStore
	Environment() EnvironmentStore
	Group() GroupStore
	HostInstance() HostInstanceStore
	Organization() OrganizationStore
	Page() PageStore
	Session() SessionStore
	User() UserStore
}

type DB interface {
	Stores
	WithTx(ctx context.Context, fn func(tx Tx) error) error
}

type Tx interface {
	Stores
}
