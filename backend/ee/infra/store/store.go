package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	apikeyStore "github.com/trysourcetool/sourcetool/backend/ee/apikey/store"
	environmentStore "github.com/trysourcetool/sourcetool/backend/ee/environment/store"
	groupStore "github.com/trysourcetool/sourcetool/backend/ee/group/store"
	hostinstanceStore "github.com/trysourcetool/sourcetool/backend/ee/hostinstance/store"
	organizationStore "github.com/trysourcetool/sourcetool/backend/ee/organization/store"
	pageStore "github.com/trysourcetool/sourcetool/backend/ee/page/store"
	sessionStore "github.com/trysourcetool/sourcetool/backend/ee/session/store"
	userStore "github.com/trysourcetool/sourcetool/backend/ee/user/store"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/session"
	"github.com/trysourcetool/sourcetool/backend/user"
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

func (s *storeEE) APIKey() apikey.Store {
	return apikeyStore.NewAPIKeyStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Environment() environment.Store {
	return environmentStore.NewEnvironmentStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Group() group.Store {
	return groupStore.NewGroupStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) HostInstance() hostinstance.Store {
	return hostinstanceStore.NewHostInstanceStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Organization() organization.Store {
	return organizationStore.NewOrganizationStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Page() page.Store {
	return pageStore.NewPageStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) Session() session.Store {
	return sessionStore.NewSessionStoreEE(infra.NewQueryLogger(s.db))
}

func (s *storeEE) User() user.Store {
	return userStore.NewUserStoreEE(infra.NewQueryLogger(s.db))
}

type transactionEE struct {
	db *sqlx.Tx
}

func (t *transactionEE) APIKey() apikey.Store {
	return apikeyStore.NewAPIKeyStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Environment() environment.Store {
	return environmentStore.NewEnvironmentStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Group() group.Store {
	return groupStore.NewGroupStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) HostInstance() hostinstance.Store {
	return hostinstanceStore.NewHostInstanceStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Organization() organization.Store {
	return organizationStore.NewOrganizationStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Page() page.Store {
	return pageStore.NewPageStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) Session() session.Store {
	return sessionStore.NewSessionStoreEE(infra.NewQueryLogger(t.db))
}

func (t *transactionEE) User() user.Store {
	return userStore.NewUserStoreEE(infra.NewQueryLogger(t.db))
}
