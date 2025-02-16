package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Organization struct {
	ID        uuid.UUID `db:"id"`
	Subdomain string    `db:"subdomain"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type (
	OrganizationByID        uuid.UUID
	OrganizationBySubdomain string
	OrganizationByUserID    uuid.UUID
)

type OrganizationStoreCE interface {
	Get(context.Context, ...any) (*Organization, error)
	Create(context.Context, *Organization) error
	IsSubdomainExists(context.Context, string) (bool, error)
}
