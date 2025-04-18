package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/config"
)

const (
	maxIdleConns = 25
	maxOpenConns = 100
)

func New() (*sqlx.DB, error) {
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
