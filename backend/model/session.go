package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	UserID         uuid.UUID `db:"user_id"`
	PageID         uuid.UUID `db:"page_id"`
	HostInstanceID uuid.UUID `db:"host_instance_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type (
	SessionByID uuid.UUID
)

type SessionStoreCE interface {
	Get(context.Context, ...any) (*Session, error)
	Create(context.Context, *Session) error
	Delete(context.Context, *Session) error
}
