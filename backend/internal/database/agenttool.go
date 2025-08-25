package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type AgentToolQuery interface{ isAgentToolQuery() }

type AgentToolByIDQuery struct{ ID uuid.UUID }

func (q AgentToolByIDQuery) isAgentToolQuery() {}

func AgentToolByID(id uuid.UUID) AgentToolQuery { return AgentToolByIDQuery{ID: id} }

type AgentToolByAgentIDQuery struct{ AgentID uuid.UUID }

func (q AgentToolByAgentIDQuery) isAgentToolQuery() {}

func AgentToolByAgentID(id uuid.UUID) AgentToolQuery {
	return AgentToolByAgentIDQuery{AgentID: id}
}

type AgentToolLimitQuery struct{ Limit uint64 }

func (q AgentToolLimitQuery) isAgentToolQuery() {}

func AgentToolLimit(limit uint64) AgentToolQuery { return AgentToolLimitQuery{Limit: limit} }

type AgentToolOffsetQuery struct{ Offset uint64 }

func (q AgentToolOffsetQuery) isAgentToolQuery() {}

func AgentToolOffset(offset uint64) AgentToolQuery { return AgentToolOffsetQuery{Offset: offset} }

type AgentToolOrderByQuery struct{ OrderBy string }

func (q AgentToolOrderByQuery) isAgentToolQuery() {}

func AgentToolOrderBy(orderBy string) AgentToolQuery { return AgentToolOrderByQuery{OrderBy: orderBy} }

type AgentToolStore interface {
	Get(ctx context.Context, queries ...AgentToolQuery) (*core.AgentTool, error)
	List(ctx context.Context, queries ...AgentToolQuery) ([]*core.AgentTool, error)
	BulkInsert(ctx context.Context, m []*core.AgentTool) error
	BulkUpdate(ctx context.Context, m []*core.AgentTool) error
	BulkDelete(ctx context.Context, m []*core.AgentTool) error
}
