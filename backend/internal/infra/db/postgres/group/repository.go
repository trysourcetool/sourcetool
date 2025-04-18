package group

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
)

type RepositoryCE struct {
	db      db.DB
	builder sq.StatementBuilderType
}

func NewRepositoryCE(db db.DB) *RepositoryCE {
	return &RepositoryCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RepositoryCE) Get(ctx context.Context, opts ...group.RepositoryOption) (*group.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) List(ctx context.Context, opts ...group.RepositoryOption) ([]*group.Group, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) Create(ctx context.Context, m *group.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) Update(ctx context.Context, m *group.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) Delete(ctx context.Context, m *group.Group) error {
	return errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) IsSlugExistsInOrganization(ctx context.Context, orgID uuid.UUID, slug string) (bool, error) {
	return false, errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) ListPages(ctx context.Context, opts ...group.PageRepositoryOption) ([]*group.GroupPage, error) {
	return nil, errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) BulkInsertPages(ctx context.Context, pages []*group.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) BulkUpdatePages(ctx context.Context, pages []*group.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}

func (r *RepositoryCE) BulkDeletePages(ctx context.Context, pages []*group.GroupPage) error {
	return errors.New("group functionality is not available in CE version")
}
