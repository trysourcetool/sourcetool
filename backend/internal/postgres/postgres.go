package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
)

var _ database.DB = (*db)(nil)

type db struct {
	db *sqlx.DB
}

func (db *db) APIKey() database.APIKeyStore {
	return newAPIKeyStore(internal.NewQueryLogger(db.db))
}

func (db *db) Environment() database.EnvironmentStore {
	return newEnvironmentStore(internal.NewQueryLogger(db.db))
}

func (db *db) Group() database.GroupStore {
	return newGroupStore(internal.NewQueryLogger(db.db))
}

func (db *db) HostInstance() database.HostInstanceStore {
	return newHostInstanceStore(internal.NewQueryLogger(db.db))
}

func (db *db) Organization() database.OrganizationStore {
	return newOrganizationStore(internal.NewQueryLogger(db.db))
}

func (db *db) Page() database.PageStore {
	return newPageStore(internal.NewQueryLogger(db.db))
}

func (db *db) Session() database.SessionStore {
	return newSessionStore(internal.NewQueryLogger(db.db))
}

func (db *db) User() database.UserStore {
	return newUserStore(internal.NewQueryLogger(db.db))
}

func (db *db) WithTx(ctx context.Context, fn func(tx database.Tx) error) error {
	return db.WithTxOptions(ctx, nil, fn)
}

func (db *db) WithTxOptions(ctx context.Context, opts *sql.TxOptions, fn func(tx database.Tx) error) error {
	if opts == nil {
		opts = &sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		}
	}

	sqlxTx, err := db.db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	t := &tx{db: sqlxTx}
	if err := fn(t); err != nil {
		if err := sqlxTx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return err
	}

	return sqlxTx.Commit()
}

func New(sqlxDB *sqlx.DB) database.DB {
	return &db{
		db: sqlxDB,
	}
}

var _ database.Tx = (*tx)(nil)

type tx struct {
	db *sqlx.Tx
}

func (tx *tx) APIKey() database.APIKeyStore {
	return newAPIKeyStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) Environment() database.EnvironmentStore {
	return newEnvironmentStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) Group() database.GroupStore {
	return newGroupStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) HostInstance() database.HostInstanceStore {
	return newHostInstanceStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) Organization() database.OrganizationStore {
	return newOrganizationStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) Page() database.PageStore {
	return newPageStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) Session() database.SessionStore {
	return newSessionStore(internal.NewQueryLogger(tx.db))
}

func (tx *tx) User() database.UserStore {
	return newUserStore(internal.NewQueryLogger(tx.db))
}

const (
	maxIdleConns = 25
	maxOpenConns = 100
)

func Open() (*sqlx.DB, error) {
	sqlDB, err := sql.Open("postgres", dsn())
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	for {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return sqlx.NewDb(sqlDB, "postgres"), nil
}

func dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.Postgres.Host,
		config.Config.Postgres.Port,
		config.Config.Postgres.User,
		config.Config.Postgres.Password,
		config.Config.Postgres.DB,
	)
}

func Migrate(dir string) error {
	db, err := Open()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+dir, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Logger.Info("Migrations applied successfully")

	return nil
}
