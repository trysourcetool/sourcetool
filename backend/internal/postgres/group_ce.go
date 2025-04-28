//go:build !ee
// +build !ee

package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
)

var _ database.GroupStore = (*groupStore)(nil)

type groupStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newGroupStore(db internal.DB) *groupStore {
	return &groupStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *groupStore) Get(ctx context.Context, queries ...database.GroupQuery) (*core.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *groupStore) List(ctx context.Context, queries ...database.GroupQuery) ([]*core.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *groupStore) Create(ctx context.Context, m *core.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *groupStore) Update(ctx context.Context, m *core.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *groupStore) Delete(ctx context.Context, m *core.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *groupStore) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	return false, errors.New("group functionality is not available in CE version")
}

func (s *groupStore) ListPages(ctx context.Context, queries ...database.GroupPageQuery) ([]*core.GroupPage, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *groupStore) BulkInsertPages(ctx context.Context, pages []*core.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *groupStore) BulkUpdatePages(ctx context.Context, pages []*core.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *groupStore) BulkDeletePages(ctx context.Context, pages []*core.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}
