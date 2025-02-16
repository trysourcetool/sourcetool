package session

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
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

func (s *StoreCE) Get(ctx context.Context, conditions ...any) (*model.Session, error) {
	query, args, err := s.buildQuery(ctx, conditions...)
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

func (s *StoreCE) buildQuery(ctx context.Context, conditions ...any) (string, []any, error) {
	q := s.builder.Select(
		`s."id"`,
		`s."organization_id"`,
		`s."user_id"`,
		`s."page_id"`,
		`s."host_instance_id"`,
		`s."created_at"`,
		`s."updated_at"`,
	).
		From(`"session" s`)

	opts, err := s.toSelectOptions(ctx, conditions...)
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	for _, o := range opts {
		q = o(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *StoreCE) toSelectOptions(ctx context.Context, conditions ...any) ([]infra.SelectOption, error) {
	options := make([]infra.SelectOption, len(conditions))
	for i, c := range conditions {
		switch v := c.(type) {
		case model.SessionByID:
			options[i] = s.byID(v)
		default:
			return nil, errdefs.ErrDatabase(errors.New("unsupported condition"))
		}
	}

	return options, nil
}

func (s *StoreCE) byID(in model.SessionByID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`s."id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) Create(ctx context.Context, m *model.Session) error {
	if _, err := s.builder.
		Insert(`"session"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"user_id"`,
			`"page_id"`,
			`"host_instance_id"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.UserID,
			m.PageID,
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
