package apikey

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/lib/pq"

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

func (s *StoreCE) Get(ctx context.Context, conditions ...any) (*model.APIKey, error) {
	query, args, err := s.buildQuery(ctx, conditions...)
	if err != nil {
		return nil, err
	}

	m := model.APIKey{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrAPIKeyNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) List(ctx context.Context, conditions ...any) ([]*model.APIKey, error) {
	query, args, err := s.buildQuery(ctx, conditions...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.APIKey, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, conditions ...any) (string, []any, error) {
	q := s.builder.Select(
		`ak."id"`,
		`ak."organization_id"`,
		`ak."environment_id"`,
		`ak."user_id"`,
		`ak."name"`,
		`ak."key"`,
		`ak."created_at"`,
		`ak."updated_at"`,
	).
		From(`"api_key" ak`)

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
		case model.APIKeyByID:
			options[i] = s.byID(v)
		case model.APIKeyByOrganizationID:
			options[i] = s.byOrganizationID(v)
		case model.APIKeyByEnvironmentID:
			options[i] = s.byEnvironmentID(v)
		case model.APIKeyByEnvironmentIDs:
			options[i] = s.byEnvironmentIDs(v)
		case model.APIKeyByUserID:
			options[i] = s.byUserID(v)
		case model.APIKeyByKey:
			options[i] = s.byKey(v)
		default:
			return nil, errdefs.ErrDatabase(errors.New("unsupported condition"))
		}
	}

	return options, nil
}

func (s *StoreCE) byID(in model.APIKeyByID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`ak."id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) byOrganizationID(in model.APIKeyByOrganizationID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`ak."organization_id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) byEnvironmentID(in model.APIKeyByEnvironmentID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`ak."environment_id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) byEnvironmentIDs(in model.APIKeyByEnvironmentIDs) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`ak."environment_id"`: []uuid.UUID(in)})
	}
}

func (s *StoreCE) byUserID(in model.APIKeyByUserID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`ak."user_id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) byKey(in model.APIKeyByKey) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`ak."key"`: in})
	}
}

func (s *StoreCE) Create(ctx context.Context, m *model.APIKey) error {
	if _, err := s.builder.
		Insert(`"api_key"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"environment_id"`,
			`"user_id"`,
			`"name"`,
			`"key"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.EnvironmentID,
			m.UserID,
			m.Name,
			m.Key,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) Update(ctx context.Context, m *model.APIKey) error {
	if _, err := s.builder.
		Update(`"api_key"`).
		Set(`"user_id"`, m.UserID).
		Set(`"name"`, m.Name).
		Set(`"key"`, m.Key).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) Delete(ctx context.Context, m *model.APIKey) error {
	if _, err := s.builder.
		Delete(`"api_key"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
