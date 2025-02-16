package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/lib/pq"
)

type Page struct {
	ID             uuid.UUID     `db:"id"`
	OrganizationID uuid.UUID     `db:"organization_id"`
	EnvironmentID  uuid.UUID     `db:"environment_id"`
	APIKeyID       uuid.UUID     `db:"api_key_id"`
	Name           string        `db:"name"`
	Route          string        `db:"route"`
	Path           pq.Int32Array `db:"path"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
}

type (
	PageByID             uuid.UUID
	PageByOrganizationID uuid.UUID
	PageByAPIKeyID       uuid.UUID
	PageBySessionID      uuid.UUID
)

type PageStoreCE interface {
	Get(context.Context, ...any) (*Page, error)
	List(context.Context, ...any) ([]*Page, error)
	BulkInsert(context.Context, []*Page) error
	BulkUpdate(context.Context, []*Page) error
	BulkDelete(context.Context, []*Page) error
}
