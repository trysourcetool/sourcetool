package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.HostInstanceStore = (*hostInstanceStore)(nil)

type hostInstanceStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newHostInstanceStore(db internal.DB) *hostInstanceStore {
	return &hostInstanceStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *hostInstanceStore) Get(ctx context.Context, queries ...database.HostInstanceQuery) (*core.HostInstance, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.HostInstance{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrHostInstanceNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *hostInstanceStore) List(ctx context.Context, queries ...database.HostInstanceQuery) ([]*core.HostInstance, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.HostInstance, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *hostInstanceStore) applyQueries(b sq.SelectBuilder, queries ...database.HostInstanceQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.HostInstanceByIDQuery:
			b = b.Where(sq.Eq{`hi."id"`: q.ID})
		case database.HostInstanceByOrganizationIDQuery:
			b = b.Where(sq.Eq{`hi."organization_id"`: q.OrganizationID})
		case database.HostInstanceByAPIKeyIDQuery:
			b = b.Where(sq.Eq{`hi."api_key_id"`: q.APIKeyID})
		case database.HostInstanceByAPIKeyQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."id" = hi."api_key_id"`).
				Where(sq.Eq{`ak."key"`: q.APIKey})
		}
	}

	return b
}

func (s *hostInstanceStore) buildQuery(ctx context.Context, queries ...database.HostInstanceQuery) (string, []any, error) {
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

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *hostInstanceStore) Create(ctx context.Context, m *core.HostInstance) error {
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

func (s *hostInstanceStore) Update(ctx context.Context, m *core.HostInstance) error {
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
