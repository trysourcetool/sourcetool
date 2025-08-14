package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Agent struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	EnvironmentID  uuid.UUID `db:"environment_id"`
	APIKeyID       uuid.UUID `db:"api_key_id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	Instructions   string    `db:"instructions"`
	Model          string    `db:"model"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type AgentTool struct {
	ID          uuid.UUID `db:"id"`
	AgentID     uuid.UUID `db:"agent_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
