//go:build !ee
// +build !ee

package postgres

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

func (db *DB) GetGroup(ctx context.Context, queries ...GroupQuery) (*core.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (db *DB) ListGroups(ctx context.Context, queries ...GroupQuery) ([]*core.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (db *DB) CreateGroup(ctx context.Context, tx *sqlx.Tx, m *core.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (db *DB) UpdateGroup(ctx context.Context, tx *sqlx.Tx, m *core.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (db *DB) DeleteGroup(ctx context.Context, tx *sqlx.Tx, m *core.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (db *DB) IsGroupSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	return false, errors.New("group functionality is not available in CE version")
}

func (db *DB) ListGroupPages(ctx context.Context, queries ...GroupPageQuery) ([]*core.GroupPage, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (db *DB) BulkInsertGroupPages(ctx context.Context, pages []*core.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (db *DB) BulkUpdateGroupPages(ctx context.Context, pages []*core.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (db *DB) BulkDeleteGroupPages(ctx context.Context, pages []*core.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}
