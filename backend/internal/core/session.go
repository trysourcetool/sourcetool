package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type SessionType int

const (
	SessionTypeUnknown SessionType = iota
	SessionTypePage
	SessionTypeAgent

	sessionTypeUnknown = "unknown"
	sessionTypePage    = "page"
	sessionTypeAgent   = "agent"
)

func (t SessionType) String() string {
	types := []string{
		sessionTypeUnknown,
		sessionTypePage,
		sessionTypeAgent,
	}

	if int(t) < 0 || int(t) >= len(types) {
		return sessionTypeUnknown
	}

	return types[t]
}

func SessionTypeFromString(s string) SessionType {
	typeMap := map[string]SessionType{
		sessionTypePage:  SessionTypePage,
		sessionTypeAgent: SessionTypeAgent,
	}

	if sessionType, ok := typeMap[s]; ok {
		return sessionType
	}
	return SessionTypeUnknown
}

type Session struct {
	ID             uuid.UUID   `db:"id"`
	OrganizationID uuid.UUID   `db:"organization_id"`
	UserID         uuid.UUID   `db:"user_id"`
	EnvironmentID  uuid.UUID   `db:"environment_id"`
	Type           SessionType `db:"type"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
}

type SessionHostInstance struct {
	ID             uuid.UUID `db:"id"`
	SessionID      uuid.UUID `db:"session_id"`
	HostInstanceID uuid.UUID `db:"host_instance_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
