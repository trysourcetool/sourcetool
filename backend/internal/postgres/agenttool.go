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

var _ database.AgentToolStore = (*agentToolStore)(nil)

type agentToolStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newAgentToolStore(db internal.DB) *agentToolStore {
	return &agentToolStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *agentToolStore) Get(ctx context.Context, queries ...database.AgentToolQuery) (*core.AgentTool, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.AgentTool{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrAgentToolNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *agentToolStore) List(ctx context.Context, queries ...database.AgentToolQuery) ([]*core.AgentTool, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.AgentTool, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *agentToolStore) applyQueries(b sq.SelectBuilder, queries ...database.AgentToolQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.AgentToolByIDQuery:
			b = b.Where(sq.Eq{`at."id"`: q.ID})
		case database.AgentToolByAgentIDQuery:
			b = b.Where(sq.Eq{`at."agent_id"`: q.AgentID})
		case database.AgentToolLimitQuery:
			b = b.Limit(q.Limit)
		case database.AgentToolOffsetQuery:
			b = b.Offset(q.Offset)
		case database.AgentToolOrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}

	return b
}

func (s *agentToolStore) buildQuery(ctx context.Context, queries ...database.AgentToolQuery) (string, []any, error) {
	q := s.builder.Select(
		`at."id"`,
		`at."agent_id"`,
		`at."name"`,
		`at."description"`,
		`at."created_at"`,
		`at."updated_at"`,
	).
		From(`"agent_tool" at`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *agentToolStore) BulkInsert(ctx context.Context, m []*core.AgentTool) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"agent_tool"`).
		Columns(
			`"id"`,
			`"agent_id"`,
			`"name"`,
			`"description"`,
		)

	for _, v := range m {
		q = q.Values(
			v.ID,
			v.AgentID,
			v.Name,
			v.Description,
		)
	}

	if _, err := q.
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *agentToolStore) BulkUpdate(ctx context.Context, m []*core.AgentTool) error {
	if len(m) == 0 {
		return nil
	}

	for _, v := range m {
		if _, err := s.builder.
			Update(`"agent_tool"`).
			Set(`"name"`, v.Name).
			Set(`"description"`, v.Description).
			Where(sq.Eq{`"id"`: v.ID}).
			RunWith(s.db).
			ExecContext(ctx); err != nil {
			return errdefs.ErrDatabase(err)
		}
	}

	return nil
}

func (s *agentToolStore) BulkDelete(ctx context.Context, m []*core.AgentTool) error {
	if len(m) == 0 {
		return nil
	}

	ids := make([]uuid.UUID, len(m))
	for i, agentTool := range m {
		ids[i] = agentTool.ID
	}

	if _, err := s.builder.
		Delete(`"agent_tool"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}