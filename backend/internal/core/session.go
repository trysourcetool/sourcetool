package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	UserID         uuid.UUID `db:"user_id"`
	EnvironmentID  uuid.UUID `db:"environment_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type SessionHostInstance struct {
	ID             uuid.UUID `db:"id"`
	SessionID      uuid.UUID `db:"session_id"`
	HostInstanceID uuid.UUID `db:"host_instance_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
