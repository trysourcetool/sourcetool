package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
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
