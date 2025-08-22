package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

type agentResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
	Model        string `json:"model"`
}

type listAgentsResponse struct {
	Agents []agentResponse `json:"agents"`
}

func (s *Server) handleListAgents(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxOrg := internal.ContextOrganization(ctx)

	environmentID := r.URL.Query().Get("environmentId")
	if environmentID == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentId is required"))
	}

	envUUID, err := uuid.FromString(environmentID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Verify environment belongs to organization
	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(envUUID))
	if err != nil {
		return err
	}
	if env.OrganizationID != ctxOrg.ID {
		return errdefs.ErrPermissionDenied(nil)
	}

	// Get agents for environment
	agents, err := s.db.Agent().List(ctx, database.AgentByEnvironmentID(envUUID))
	if err != nil {
		return err
	}

	response := listAgentsResponse{
		Agents: make([]agentResponse, len(agents)),
	}

	for i, agent := range agents {
		response.Agents[i] = agentResponse{
			ID:           agent.ID.String(),
			Name:         agent.Name,
			Description:  agent.Description,
			Instructions: agent.Instructions,
			Model:        agent.Model,
		}
	}

	return s.renderJSON(w, http.StatusOK, response)
}

func (s *Server) handleGetAgent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxOrg := internal.ContextOrganization(ctx)

	agentIDStr := chi.URLParam(r, "agentID")
	agentID, err := uuid.FromString(agentIDStr)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	agent, err := s.db.Agent().Get(ctx, database.AgentByID(agentID))
	if err != nil {
		return err
	}

	// Verify agent belongs to organization
	if agent.OrganizationID != ctxOrg.ID {
		return errdefs.ErrPermissionDenied(nil)
	}

	response := agentResponse{
		ID:           agent.ID.String(),
		Name:         agent.Name,
		Description:  agent.Description,
		Instructions: agent.Instructions,
		Model:        agent.Model,
	}

	return s.renderJSON(w, http.StatusOK, response)
}