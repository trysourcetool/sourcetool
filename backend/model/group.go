package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/storeopts"
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

type GroupStore interface {
	Get(context.Context, ...storeopts.GroupOption) (*Group, error)
	List(context.Context, ...storeopts.GroupOption) ([]*Group, error)
	Create(context.Context, *Group) error
	Update(context.Context, *Group) error
	Delete(context.Context, *Group) error
	IsSlugExistsInOrganization(context.Context, uuid.UUID, string) (bool, error)

	ListPages(context.Context, ...storeopts.GroupPageOption) ([]*GroupPage, error)
	BulkInsertPages(context.Context, []*GroupPage) error
	BulkUpdatePages(context.Context, []*GroupPage) error
	BulkDeletePages(context.Context, []*GroupPage) error
}
