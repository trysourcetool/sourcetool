package store

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type GroupStoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewGroupStoreCE(db infra.DB) *GroupStoreCE {
	return &GroupStoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *GroupStoreCE) Get(ctx context.Context, opts ...group.StoreOption) (*group.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) List(ctx context.Context, opts ...group.StoreOption) ([]*group.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) Create(ctx context.Context, m *group.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) Update(ctx context.Context, m *group.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) Delete(ctx context.Context, m *group.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	return false, errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) ListPages(ctx context.Context, opts ...group.PageStoreOption) ([]*group.GroupPage, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) BulkInsertPages(ctx context.Context, pages []*group.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) BulkUpdatePages(ctx context.Context, pages []*group.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *GroupStoreCE) BulkDeletePages(ctx context.Context, pages []*group.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}
