package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type HostInstanceStoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewHostInstanceStoreCE(db infra.DB) *HostInstanceStoreCE {
	return &HostInstanceStoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *HostInstanceStoreCE) Get(ctx context.Context, opts ...hostinstance.StoreOption) (*hostinstance.HostInstance, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := hostinstance.HostInstance{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrHostInstanceNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *HostInstanceStoreCE) List(ctx context.Context, opts ...hostinstance.StoreOption) ([]*hostinstance.HostInstance, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*hostinstance.HostInstance, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *HostInstanceStoreCE) buildQuery(ctx context.Context, opts ...hostinstance.StoreOption) (string, []any, error) {
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

func (s *HostInstanceStoreCE) Create(ctx context.Context, m *hostinstance.HostInstance) error {
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

func (s *HostInstanceStoreCE) Update(ctx context.Context, m *hostinstance.HostInstance) error {
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
