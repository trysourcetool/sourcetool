package infra

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type ModelStore interface {
	APIKey() model.APIKeyStore
	Environment() model.EnvironmentStore
	Group() model.GroupStore
	HostInstance() model.HostInstanceStore
	Organization() model.OrganizationStore
	Page() model.PageStore
	Session() model.SessionStore
	User() model.UserStore
}

type Store interface {
	ModelStore
	Close() error
	RunTransaction(func(tx Transaction) error) error
}

type Transaction interface {
	ModelStore
}

type DB interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	Exec(string, ...any) (sql.Result, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	GetContext(context.Context, any, string, ...any) error
	QueryxContext(context.Context, string, ...any) (*sqlx.Rows, error)
	SelectContext(context.Context, any, string, ...any) error
}

type queryLogger struct {
	db DB
}

func NewQueryLogger(db DB) *queryLogger {
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

type (
	SelectOption      func(b sq.SelectBuilder) sq.SelectBuilder
	LoadOption[T any] func(ctx context.Context, e ...T) error
	Limit             uint64
	Offset            uint64
	OrderBy           string
	GroupBy           string
)
