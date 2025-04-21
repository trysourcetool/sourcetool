package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
)

type db interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	Exec(string, ...any) (sql.Result, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	GetContext(context.Context, any, string, ...any) error
	QueryxContext(context.Context, string, ...any) (*sqlx.Rows, error)
	SelectContext(context.Context, any, string, ...any) error
	Beginx() (*sqlx.Tx, error)
}

type queryLogger struct {
	db db
}

func NewQueryLogger(db db) *queryLogger {
	return &queryLogger{db}
}

func (l *queryLogger) Query(query string, args ...any) (*sql.Rows, error) {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.Query(query, args...)
}

func (l *queryLogger) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.QueryContext(ctx, query, args...)
}

func (l *queryLogger) Exec(query string, args ...any) (sql.Result, error) {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.Exec(query, args...)
}

func (l *queryLogger) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.ExecContext(ctx, query, args...)
}

func (l *queryLogger) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.GetContext(ctx, dest, query, args...)
}

func (l *queryLogger) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.QueryxContext(ctx, query, args...)
}

func (l *queryLogger) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	logger.Logger.Sugar().Debugf("%s, args: %s", query, args)
	return l.db.SelectContext(ctx, dest, query, args...)
}

func (l *queryLogger) Beginx() (*sqlx.Tx, error) {
	return l.db.Beginx()
}

type DB struct {
	db      db
	builder sq.StatementBuilderType
}

func New(db db) *DB {
	return &DB{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (db *DB) Beginx() (*sqlx.Tx, error) {
	return db.db.Beginx()
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

	return nil
}
