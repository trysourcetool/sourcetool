package session

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

func (s *StoreCE) Get(ctx context.Context, opts ...storeopts.SessionOption) (*model.Session, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.Session{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrSessionNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, opts ...storeopts.SessionOption) (string, []any, error) {
	q := s.builder.Select(
		`s."id"`,
		`s."organization_id"`,
		`s."user_id"`,
		`s."api_key_id"`,
		`s."host_instance_id"`,
		`s."created_at"`,
		`s."updated_at"`,
	).
		From(`"session" s`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *StoreCE) Create(ctx context.Context, m *model.Session) error {
	if _, err := s.builder.
		Insert(`"session"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"user_id"`,
			`"api_key_id"`,
			`"host_instance_id"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.UserID,
			m.APIKeyID,
			m.HostInstanceID,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) Delete(ctx context.Context, m *model.Session) error {
	if _, err := s.builder.
		Delete(`"session"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
