package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Group struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Name           string    `db:"name"`
	Slug           string    `db:"slug"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type GroupPage struct {
	ID        uuid.UUID `db:"id"`
	GroupID   uuid.UUID `db:"group_id"`
	PageID    uuid.UUID `db:"page_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
