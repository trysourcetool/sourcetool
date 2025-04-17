package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	apikeyStore "github.com/trysourcetool/sourcetool/backend/apikey/store"
	"github.com/trysourcetool/sourcetool/backend/environment"
	environmentStore "github.com/trysourcetool/sourcetool/backend/environment/store"
	"github.com/trysourcetool/sourcetool/backend/group"
	groupStore "github.com/trysourcetool/sourcetool/backend/group/store"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	hostinstanceStore "github.com/trysourcetool/sourcetool/backend/hostinstance/store"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization"
	organizationStore "github.com/trysourcetool/sourcetool/backend/organization/store"
	"github.com/trysourcetool/sourcetool/backend/page"
	pageStore "github.com/trysourcetool/sourcetool/backend/page/store"
	"github.com/trysourcetool/sourcetool/backend/session"
	sessionStore "github.com/trysourcetool/sourcetool/backend/session/store"
	"github.com/trysourcetool/sourcetool/backend/user"
	userStore "github.com/trysourcetool/sourcetool/backend/user/store"
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

func (s *storeCE) APIKey() apikey.Store {
	return apikeyStore.NewAPIKeyStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Environment() environment.Store {
	return environmentStore.NewEnvironmentStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Group() group.Store {
	return groupStore.NewGroupStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) HostInstance() hostinstance.Store {
	return hostinstanceStore.NewHostInstanceStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Organization() organization.Store {
	return organizationStore.NewOrganizationStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Page() page.Store {
	return pageStore.NewPageStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) Session() session.Store {
	return sessionStore.NewSessionStoreCE(infra.NewQueryLogger(s.db))
}

func (s *storeCE) User() user.Store {
	return userStore.NewUserStoreCE(infra.NewQueryLogger(s.db))
}

type transactionCE struct {
	db *sqlx.Tx
}

func (t *transactionCE) APIKey() apikey.Store {
	return apikeyStore.NewAPIKeyStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Environment() environment.Store {
	return environmentStore.NewEnvironmentStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Group() group.Store {
	return groupStore.NewGroupStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) HostInstance() hostinstance.Store {
	return hostinstanceStore.NewHostInstanceStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Organization() organization.Store {
	return organizationStore.NewOrganizationStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Page() page.Store {
	return pageStore.NewPageStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) Session() session.Store {
	return sessionStore.NewSessionStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transactionCE) User() user.Store {
	return userStore.NewUserStoreCE(infra.NewQueryLogger(t.db))
}
