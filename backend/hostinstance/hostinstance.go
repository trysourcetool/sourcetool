package hostinstance

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type HostInstance struct {
	ID             uuid.UUID          `db:"id"`
	OrganizationID uuid.UUID          `db:"organization_id"`
	APIKeyID       uuid.UUID          `db:"api_key_id"`
	SDKName        string             `db:"sdk_name"`
	SDKVersion     string             `db:"sdk_version"`
	Status         HostInstanceStatus `db:"status"`
	CreatedAt      time.Time          `db:"created_at"`
	UpdatedAt      time.Time          `db:"updated_at"`
}
