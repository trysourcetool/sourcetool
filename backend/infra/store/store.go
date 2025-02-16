package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/session"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type StoreCE struct {
	db *sqlx.DB
}

func NewCE(db *sqlx.DB) *StoreCE {
	return &StoreCE{
		db: db,
	}
}

func (s *StoreCE) Close() error {
	return s.db.Close()
}

func (s *StoreCE) RunTransaction(f func(infra.Transaction) error) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	t := &transaction{db: tx}
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

func (s *StoreCE) APIKey() model.APIKeyStore {
	return apikey.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) Environment() model.EnvironmentStoreCE {
	return environment.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) Group() model.GroupStoreCE {
	return group.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) HostInstance() model.HostInstanceStoreCE {
	return hostinstance.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) Organization() model.OrganizationStoreCE {
	return organization.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) Page() model.PageStoreCE {
	return page.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) Session() model.SessionStoreCE {
	return session.NewStoreCE(infra.NewQueryLogger(s.db))
}

func (s *StoreCE) User() model.UserStoreCE {
	return user.NewStoreCE(infra.NewQueryLogger(s.db))
}

type transaction struct {
	db *sqlx.Tx
}

func (t *transaction) APIKey() model.APIKeyStore {
	return apikey.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) Environment() model.EnvironmentStoreCE {
	return environment.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) Group() model.GroupStoreCE {
	return group.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) HostInstance() model.HostInstanceStoreCE {
	return hostinstance.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) Organization() model.OrganizationStoreCE {
	return organization.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) Page() model.PageStoreCE {
	return page.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) Session() model.SessionStoreCE {
	return session.NewStoreCE(infra.NewQueryLogger(t.db))
}

func (t *transaction) User() model.UserStoreCE {
	return user.NewStoreCE(infra.NewQueryLogger(t.db))
}
