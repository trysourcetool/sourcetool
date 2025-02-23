package group

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type StoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewStoreCE(db infra.DB) *StoreCE {
	return &StoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *StoreCE) Get(ctx context.Context, conditions ...any) (*model.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) List(ctx context.Context, conditions ...any) ([]*model.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) Create(ctx context.Context, m *model.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) Update(ctx context.Context, m *model.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) Delete(ctx context.Context, m *model.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	return false, errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) ListPages(ctx context.Context, conditions ...any) ([]*model.GroupPage, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) BulkInsertPages(ctx context.Context, pages []*model.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) BulkUpdatePages(ctx context.Context, pages []*model.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (s *StoreCE) BulkDeletePages(ctx context.Context, pages []*model.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}
