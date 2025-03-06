package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/ee/apikey"
	"github.com/trysourcetool/sourcetool/backend/ee/environment"
	"github.com/trysourcetool/sourcetool/backend/ee/group"
	"github.com/trysourcetool/sourcetool/backend/ee/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/ee/organization"
	"github.com/trysourcetool/sourcetool/backend/ee/page"
	"github.com/trysourcetool/sourcetool/backend/ee/session"
	"github.com/trysourcetool/sourcetool/backend/ee/user"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type storeEE struct {
	db *sqlx.DB
}

func NewEE(db *sqlx.DB) *storeEE {
	return &storeEE{
		db: db,
	}
}

func (s *storeEE) Close() error {
	return s.db.Close()
}

func (s *storeEE) RunTransaction(f func(infra.Transaction) error) error {
	tx, err := s.db.Beginx()
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

func (s *storeEE) APIKey() model.APIKeyStore {
	return apikey.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Environment() model.EnvironmentStore {
	return environment.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Group() model.GroupStore {
	return group.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) HostInstance() model.HostInstanceStore {
	return hostinstance.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Organization() model.OrganizationStore {
	return organization.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Page() model.PageStore {
	return page.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Session() model.SessionStore {
	return session.NewStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) User() model.UserStore {
	return user.NewStoreEE(infra.NewQueryLogger(s.db))
}

type transactionEE struct {
	db *sqlx.Tx
}

func (t *transactionEE) APIKey() model.APIKeyStore {
	return apikey.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Environment() model.EnvironmentStore {
	return environment.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Group() model.GroupStore {
	return group.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) HostInstance() model.HostInstanceStore {
	return hostinstance.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Organization() model.OrganizationStore {
	return organization.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Page() model.PageStore {
	return page.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Session() model.SessionStore {
	return session.NewStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) User() model.UserStore {
	return user.NewStoreEE(infra.NewQueryLogger(t.db))
}
