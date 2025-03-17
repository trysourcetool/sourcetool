package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type APIKey struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	EnvironmentID  uuid.UUID `db:"environment_id"`
	UserID         uuid.UUID `db:"user_id"`
	Name           string    `db:"name"`
	Key            string    `db:"key"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type APIKeyStore interface {
	Get(context.Context, ...storeopts.APIKeyOption) (*APIKey, error)
	List(context.Context, ...storeopts.APIKeyOption) ([]*APIKey, error)
	Create(context.Context, *APIKey) error
	Update(context.Context, *APIKey) error
	Delete(context.Context, *APIKey) error
}
