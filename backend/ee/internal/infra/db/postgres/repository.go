package postgres

import (
	"github.com/jmoiron/sqlx"

	apikeyRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/apikey"
	environmentRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/environment"
	groupRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/group"
	hostinstanceRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/hostinstance"
	organizationRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/organization"
	pageRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/page"
	sessionRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/session"
	userRepo "github.com/trysourcetool/sourcetool/backend/ee/internal/infra/db/postgres/user"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
)

type repositoryEE struct {
	db *sqlx.DB
}

func NewRepositoryEE(db *sqlx.DB) *repositoryEE {
	return &repositoryEE{
		db: db,
	}
}

func (r *repositoryEE) Close() error {
	return r.db.Close()
}

func (r *repositoryEE) RunTransaction(f func(db.Transaction) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	t := &transactionEE{db: tx}
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

func (r *repositoryEE) APIKey() apikey.Repository {
	return apikeyRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) Environment() environment.Repository {
	return environmentRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) Group() group.Repository {
	return groupRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) HostInstance() hostinstance.Repository {
	return hostinstanceRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) Organization() organization.Repository {
	return organizationRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) Page() page.Repository {
	return pageRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) Session() session.Repository {
	return sessionRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

func (r *repositoryEE) User() user.Repository {
	return userRepo.NewRepositoryEE(db.NewQueryLogger(r.db))
}

type transactionEE struct {
	db *sqlx.Tx
}

func (t *transactionEE) APIKey() apikey.Repository {
	return apikeyRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) Environment() environment.Repository {
	return environmentRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) Group() group.Repository {
	return groupRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) HostInstance() hostinstance.Repository {
	return hostinstanceRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) Organization() organization.Repository {
	return organizationRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) Page() page.Repository {
	return pageRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) Session() session.Repository {
	return sessionRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}

func (t *transactionEE) User() user.Repository {
	return userRepo.NewRepositoryEE(db.NewQueryLogger(t.db))
}
