package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.AgentStore = (*agentStore)(nil)

type agentStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newAgentStore(db internal.DB) *agentStore {
	return &agentStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *agentStore) Get(ctx context.Context, queries ...database.AgentQuery) (*core.Agent, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Agent{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrAgentNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *agentStore) List(ctx context.Context, queries ...database.AgentQuery) ([]*core.Agent, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.Agent, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *agentStore) applyQueries(b sq.SelectBuilder, queries ...database.AgentQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.AgentByIDQuery:
			b = b.Where(sq.Eq{`a."id"`: q.ID})
		case database.AgentByOrganizationIDQuery:
			b = b.Where(sq.Eq{`a."organization_id"`: q.OrganizationID})
		case database.AgentByAPIKeyIDQuery:
			b = b.Where(sq.Eq{`a."api_key_id"`: q.APIKeyID})
		case database.AgentByEnvironmentIDQuery:
			b = b.Where(sq.Eq{`a."environment_id"`: q.EnvironmentID})
		case database.AgentBySessionIDQuery:
			b = b.
				InnerJoin(`"api_key" ak ON ak."id" = a."api_key_id"`).
				InnerJoin(`"environment" e ON e."id" = ak."environment_id"`).
				InnerJoin(`"session" s ON s."environment_id" = e."id"`).
				Where(sq.Eq{`s."id"`: q.SessionID})
		case database.AgentLimitQuery:
			b = b.Limit(q.Limit)
		case database.AgentOffsetQuery:
			b = b.Offset(q.Offset)
		case database.AgentOrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}

	return b
}

func (s *agentStore) buildQuery(ctx context.Context, queries ...database.AgentQuery) (string, []any, error) {
	q := s.builder.Select(
		`a."id"`,
		`a."organization_id"`,
		`a."environment_id"`,
		`a."api_key_id"`,
		`a."name"`,
		`a."description"`,
		`a."instructions"`,
		`a."model"`,
		`a."created_at"`,
		`a."updated_at"`,
	).
		From(`"agent" a`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *agentStore) BulkInsert(ctx context.Context, m []*core.Agent) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"agent"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"environment_id"`,
			`"api_key_id"`,
			`"name"`,
			`"description"`,
			`"instructions"`,
			`"model"`,
		)

	for _, v := range m {
		q = q.Values(
			v.ID,
			v.OrganizationID,
			v.EnvironmentID,
			v.APIKeyID,
			v.Name,
			v.Description,
			v.Instructions,
			v.Model,
		)
	}

	if _, err := q.
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *agentStore) BulkUpdate(ctx context.Context, m []*core.Agent) error {
	if len(m) == 0 {
		return nil
	}

	for _, v := range m {
		if _, err := s.builder.
			Update(`"agent"`).
			Set(`"name"`, v.Name).
			Set(`"description"`, v.Description).
			Set(`"instructions"`, v.Instructions).
			Set(`"model"`, v.Model).
			Where(sq.Eq{`"id"`: v.ID}).
			RunWith(s.db).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (s *agentStore) BulkDelete(ctx context.Context, m []*core.Agent) error {
	if len(m) == 0 {
		return nil
	}

	ids := make([]uuid.UUID, len(m))
	for i, agent := range m {
		ids[i] = agent.ID
	}

	if _, err := s.builder.
		Delete(`"agent"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
