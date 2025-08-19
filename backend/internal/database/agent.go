package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type AgentQuery interface{ isAgentQuery() }

type AgentByIDQuery struct{ ID uuid.UUID }

func (q AgentByIDQuery) isAgentQuery() {}

func AgentByID(id uuid.UUID) AgentQuery { return AgentByIDQuery{ID: id} }

type AgentByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q AgentByOrganizationIDQuery) isAgentQuery() {}

func AgentByOrganizationID(id uuid.UUID) AgentQuery {
	return AgentByOrganizationIDQuery{OrganizationID: id}
}

type AgentByAPIKeyIDQuery struct{ APIKeyID uuid.UUID }

func (q AgentByAPIKeyIDQuery) isAgentQuery() {}

func AgentByAPIKeyID(id uuid.UUID) AgentQuery { return AgentByAPIKeyIDQuery{APIKeyID: id} }

type AgentByEnvironmentIDQuery struct{ EnvironmentID uuid.UUID }

func (q AgentByEnvironmentIDQuery) isAgentQuery() {}

func AgentByEnvironmentID(id uuid.UUID) AgentQuery {
	return AgentByEnvironmentIDQuery{EnvironmentID: id}
}

type AgentLimitQuery struct{ Limit uint64 }

func (q AgentLimitQuery) isAgentQuery() {}

func AgentLimit(limit uint64) AgentQuery { return AgentLimitQuery{Limit: limit} }

type AgentOffsetQuery struct{ Offset uint64 }

func (q AgentOffsetQuery) isAgentQuery() {}

func AgentOffset(offset uint64) AgentQuery { return AgentOffsetQuery{Offset: offset} }

type AgentOrderByQuery struct{ OrderBy string }

func (q AgentOrderByQuery) isAgentQuery() {}

func AgentOrderBy(orderBy string) AgentQuery { return AgentOrderByQuery{OrderBy: orderBy} }

type AgentStore interface {
	Get(ctx context.Context, queries ...AgentQuery) (*core.Agent, error)
	List(ctx context.Context, queries ...AgentQuery) ([]*core.Agent, error)
	BulkInsert(ctx context.Context, m []*core.Agent) error
	BulkUpdate(ctx context.Context, m []*core.Agent) error
	BulkDelete(ctx context.Context, m []*core.Agent) error
}
