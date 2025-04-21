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
	HostInstanceStatusOffline
	HostInstanceStatusShuttingDown

	hostInstanceStatusUnknown      = "unknown"
	hostInstanceStatusOnline       = "online"
	hostInstanceStatusUnreachable  = "unreachable"
	hostInstanceStatusOffline      = "offline"
	hostInstanceStatusShuttingDown = "shuttingDown"
)

func (s HostInstanceStatus) String() string {
	switch s {
	case HostInstanceStatusOnline:
		return hostInstanceStatusOnline
	case HostInstanceStatusUnreachable:
		return hostInstanceStatusUnreachable
	case HostInstanceStatusOffline:
		return hostInstanceStatusOffline
	case HostInstanceStatusShuttingDown:
		return hostInstanceStatusShuttingDown
	default:
		return hostInstanceStatusUnknown
	}
}

func HostInstanceStatusFromString(s string) HostInstanceStatus {
	switch s {
	case hostInstanceStatusOnline:
		return HostInstanceStatusOnline
	case hostInstanceStatusUnreachable:
		return HostInstanceStatusUnreachable
	case hostInstanceStatusOffline:
		return HostInstanceStatusOffline
	case hostInstanceStatusShuttingDown:
		return HostInstanceStatusShuttingDown
	default:
		return HostInstanceStatusUnknown
	}
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
