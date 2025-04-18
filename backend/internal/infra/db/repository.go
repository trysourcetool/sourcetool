package db

import (
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/session"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
)

type Repositories interface {
	APIKey() apikey.Repository
	Environment() environment.Repository
	Group() group.Repository
	HostInstance() hostinstance.Repository
	Organization() organization.Repository
	Page() page.Repository
	Session() session.Repository
	User() user.Repository
}

type Repository interface {
	Repositories
	Close() error
	RunTransaction(func(tx Transaction) error) error
}

type Transaction interface {
	Repositories
}
