package database

import (
	"context"
	"database/sql"
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
	WithTxOptions(ctx context.Context, opts *sql.TxOptions, fn func(tx Tx) error) error
}

type Tx interface {
	Stores
}
