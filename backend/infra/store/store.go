package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/health"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/session"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type storeCE struct {
	db *sqlx.DB
}

func NewCE(db *sqlx.DB) *storeCE {
	return &storeCE{
		db: db,
	}
}

func (s *storeCE) Close() error {
	return s.db.Close()
}

func (s *storeCE) RunTransaction(f func(infra.Transaction) error) error {
	tx, err := s.db.Beginx()
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

func (s *storeCE) APIKey() model.APIKeyStore {
	return apikey.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Environment() model.EnvironmentStore {
	return environment.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Group() model.GroupStore {
	return group.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) HostInstance() model.HostInstanceStore {
	return hostinstance.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Organization() model.OrganizationStore {
	return organization.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Page() model.PageStore {
	return page.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Session() model.SessionStore {
	return session.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) User() model.UserStore {
	return user.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Health() model.HealthStore {
	return health.NewStoreCE(infra.NewQueryLogger(s.db))
}

type transactionCE struct {
	db *sqlx.Tx
}

func (t *transactionCE) APIKey() model.APIKeyStore {
	return apikey.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Environment() model.EnvironmentStore {
	return environment.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Group() model.GroupStore {
	return group.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) HostInstance() model.HostInstanceStore {
	return hostinstance.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Organization() model.OrganizationStore {
	return organization.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Page() model.PageStore {
	return page.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Session() model.SessionStore {
	return session.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) User() model.UserStore {
	return user.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Health() model.HealthStore {
	return health.NewStoreCE(infra.NewQueryLogger(t.db))
}
