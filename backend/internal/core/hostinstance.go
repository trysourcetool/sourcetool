package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type HostInstanceStatus int

const (
	HostInstanceStatusUnknown HostInstanceStatus = iota
	HostInstanceStatusOnline
	HostInstanceStatusUnreachable

	hostInstanceStatusUnknown     = "unknown"
	hostInstanceStatusOnline      = "online"
	hostInstanceStatusUnreachable = "unreachable"
)

func (s HostInstanceStatus) String() string {
	statuses := []string{
		hostInstanceStatusUnknown,
		hostInstanceStatusOnline,
		hostInstanceStatusUnreachable,
	}
	return statuses[s]
}

func HostInstanceStatusFromString(s string) HostInstanceStatus {
	statusMap := map[string]HostInstanceStatus{
		hostInstanceStatusOnline:      HostInstanceStatusOnline,
		hostInstanceStatusUnreachable: HostInstanceStatusUnreachable,
	}

	if status, ok := statusMap[s]; ok {
		return status
	}
	return HostInstanceStatusUnknown
}

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
