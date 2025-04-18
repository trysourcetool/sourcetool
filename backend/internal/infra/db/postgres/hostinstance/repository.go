package hostinstance

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal/domain/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
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

func (r *RepositoryCE) Get(ctx context.Context, opts ...hostinstance.RepositoryOption) (*hostinstance.HostInstance, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := hostinstance.HostInstance{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrHostInstanceNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) List(ctx context.Context, opts ...hostinstance.RepositoryOption) ([]*hostinstance.HostInstance, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*hostinstance.HostInstance, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildQuery(ctx context.Context, opts ...hostinstance.RepositoryOption) (string, []any, error) {
	q := r.builder.Select(
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

func (r *RepositoryCE) Create(ctx context.Context, m *hostinstance.HostInstance) error {
	if _, err := r.builder.
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
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) Update(ctx context.Context, m *hostinstance.HostInstance) error {
	if _, err := r.builder.
		Update(`"host_instance"`).
		Set(`"sdk_name"`, m.SDKName).
		Set(`"sdk_version"`, m.SDKVersion).
		Set(`"status"`, m.Status).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
