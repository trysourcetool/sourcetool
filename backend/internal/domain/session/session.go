package session

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	UserID         uuid.UUID `db:"user_id"`
	APIKeyID       uuid.UUID `db:"api_key_id"`
	HostInstanceID uuid.UUID `db:"host_instance_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
