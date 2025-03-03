package hostinstance

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
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

func (s *StoreCE) Get(ctx context.Context, opts ...storeopts.HostInstanceOption) (*model.HostInstance, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.HostInstance{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrHostInstanceNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) List(ctx context.Context, opts ...storeopts.HostInstanceOption) ([]*model.HostInstance, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.HostInstance, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, opts ...storeopts.HostInstanceOption) (string, []any, error) {
	q := s.builder.Select(
		`hi."id"`,
		`hi."organization_id"`,
		`hi."api_key_id"`,
		`hi."sdk_name"`,
		`hi."sdk_version"`,
		`hi."status"`,
		`hi."created_at"`,
		`hi."updated_at"`,
	).
		From(`"host_instance" hi`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *StoreCE) Create(ctx context.Context, m *model.HostInstance) error {
	if _, err := s.builder.
		Insert(`"host_instance"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"api_key_id"`,
			`"sdk_name"`,
			`"sdk_version"`,
			`"status"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.APIKeyID,
			m.SDKName,
			m.SDKVersion,
			m.Status,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) Update(ctx context.Context, m *model.HostInstance) error {
	if _, err := s.builder.
		Update(`"host_instance"`).
		Set(`"sdk_name"`, m.SDKName).
		Set(`"sdk_version"`, m.SDKVersion).
		Set(`"status"`, m.Status).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
