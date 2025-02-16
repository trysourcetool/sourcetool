package postgres

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"

	"github.com/trysourcetool/sourcetool/backend/logger"
)

func Migrate() error {
	migrations := &migrate.FileMigrationSource{
		Dir: "postgres/ddl",
	}

	db, err := New()
	if err != nil {
		return err
	}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	logger.Logger.Info(fmt.Sprintf("Applied %d migrations!", n))

	return nil
}
