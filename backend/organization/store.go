package organization

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
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

func (s *StoreCE) Get(ctx context.Context, conditions ...any) (*model.Organization, error) {
	query, args, err := s.buildQuery(ctx, conditions...)
	if err != nil {
		return nil, err
	}

	m := model.Organization{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrOrganizationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, conditions ...any) (string, []any, error) {
	q := s.builder.Select(
		`o."id"`,
		`o."subdomain"`,
		`o."created_at"`,
		`o."updated_at"`,
	).
		From(`"organization" o`)

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
		case model.OrganizationByID:
			options[i] = s.byID(v)
		case model.OrganizationBySubdomain:
			options[i] = s.bySubdomain(v)
		case model.OrganizationByUserID:
			options[i] = s.byUserID(v)
		case infra.Limit:
			options[i] = s.limit(v)
		case infra.Offset:
			options[i] = s.offset(v)
		case infra.OrderBy:
			options[i] = s.orderBy(v)
		default:
			return nil, errdefs.ErrDatabase(errors.New("unsupported condition"))
		}
	}

	return options, nil
}

func (s *StoreCE) byID(in model.OrganizationByID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`o."id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) bySubdomain(in model.OrganizationBySubdomain) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{`o."subdomain"`: in})
	}
}

func (s *StoreCE) byUserID(in model.OrganizationByUserID) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.
			InnerJoin(`"user_organization_access" uoa ON uoa."organization_id" = o."id"`).
			Where(sq.Eq{`uoa."user_id"`: uuid.UUID(in)})
	}
}

func (s *StoreCE) Create(ctx context.Context, m *model.Organization) error {
	if _, err := s.builder.
		Insert(`"organization"`).
		Columns(
			`"id"`,
			`"subdomain"`,
		).
		Values(
			m.ID,
			m.Subdomain,
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

func (s *StoreCE) limit(in infra.Limit) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Limit(uint64(in))
	}
}

func (s *StoreCE) offset(in infra.Offset) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Offset(uint64(in))
	}
}

func (s *StoreCE) orderBy(in infra.OrderBy) infra.SelectOption {
	return func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.OrderBy(string(in))
	}
}

func (s *StoreCE) IsSubdomainExists(ctx context.Context, subdomain string) (bool, error) {
	if _, err := s.Get(ctx, model.OrganizationBySubdomain(subdomain)); err != nil {
		if errdefs.IsOrganizationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
