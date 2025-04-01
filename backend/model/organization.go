package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type Organization struct {
	ID        uuid.UUID `db:"id"`
	Subdomain *string   `db:"subdomain"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrganizationStore interface {
	Get(context.Context, ...storeopts.OrganizationOption) (*Organization, error)
	List(context.Context, ...storeopts.OrganizationOption) ([]*Organization, error)
	Create(context.Context, *Organization) error
	IsSubdomainExists(context.Context, string) (bool, error)
}
