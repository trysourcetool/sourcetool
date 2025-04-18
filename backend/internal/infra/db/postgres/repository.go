package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	apikeyRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/apikey"
	environmentRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/environment"
	groupRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/group"
	hostinstanceRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/hostinstance"
	organizationRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/organization"
	pageRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/page"
	sessionRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/session"
	userRepo "github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres/user"
)

type repositoryCE struct {
	db *sqlx.DB
}

func NewRepositoryCE(db *sqlx.DB) *repositoryCE {
	return &repositoryCE{
		db: db,
	}
}

func (r *repositoryCE) Close() error {
	return r.db.Close()
}

func (r *repositoryCE) RunTransaction(f func(db.Transaction) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	t := &transactionCE{db: tx}
	if err := f(t); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *repositoryCE) APIKey() apikey.Repository {
	return apikeyRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) Environment() environment.Repository {
	return environmentRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) Group() group.Repository {
	return groupRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) HostInstance() hostinstance.Repository {
	return hostinstanceRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) Organization() organization.Repository {
	return organizationRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) Page() page.Repository {
	return pageRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) Session() session.Repository {
	return sessionRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

func (r *repositoryCE) User() user.Repository {
	return userRepo.NewRepositoryCE(db.NewQueryLogger(r.db))
}

type transactionCE struct {
	db *sqlx.Tx
}

func (t *transactionCE) APIKey() apikey.Repository {
	return apikeyRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) Environment() environment.Repository {
	return environmentRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) Group() group.Repository {
	return groupRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) HostInstance() hostinstance.Repository {
	return hostinstanceRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) Organization() organization.Repository {
	return organizationRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) Page() page.Repository {
	return pageRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) Session() session.Repository {
	return sessionRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}

func (t *transactionCE) User() user.Repository {
	return userRepo.NewRepositoryCE(db.NewQueryLogger(t.db))
}
