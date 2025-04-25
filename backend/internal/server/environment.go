package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

type environmentResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func environmentFromModel(env *core.Environment) *environmentResponse {
	if env == nil {
		return nil
	}

	return &environmentResponse{
		ID:        env.ID.String(),
		Name:      env.Name,
		Slug:      env.Slug,
		Color:     env.Color,
		CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
	}
}

type getEnvironmentResponse struct {
	Environment *environmentResponse `json:"environment"`
}

func (s *Server) handleGetEnvironment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	envIDReq := chi.URLParam(r, "environmentID")
	if envIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentID is required"))
	}

	envID, err := uuid.FromString(envIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	ctxOrg := internal.ContextOrganization(ctx)

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByOrganizationID(ctxOrg.ID), database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, getEnvironmentResponse{
		Environment: environmentFromModel(env),
	})
}

type listEnvironmentsResponse struct {
	Environments []*environmentResponse `json:"environments"`
}

func (s *Server) handleListEnvironments(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxOrg := internal.ContextOrganization(ctx)
	envs, err := s.db.Environment().List(ctx, database.EnvironmentByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	envsOut := make([]*environmentResponse, 0, len(envs))
	for _, env := range envs {
		envsOut = append(envsOut, environmentFromModel(env))
	}

	return s.renderJSON(w, http.StatusOK, listEnvironmentsResponse{
		Environments: envsOut,
	})
}

type createEnvironmentRequest struct {
	Name  string `json:"name" validate:"required"`
	Slug  string `json:"slug" validate:"required"`
	Color string `json:"color" validate:"required"`
}

type createEnvironmentResponse struct {
	Environment *environmentResponse `json:"environment"`
}

func (s *Server) handleCreateEnvironment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req createEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditEnvironment); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)

	slugExists, err := s.db.Environment().IsSlugExistsInOrganization(ctx, ctxOrg.ID, req.Slug)
	if err != nil {
		return err
	}
	if slugExists {
		return errdefs.ErrEnvironmentSlugAlreadyExists(errors.New("slug already exists"))
	}

	if !validateSlug(req.Slug) {
		return errdefs.ErrInvalidArgument(errors.New("invalid slug"))
	}

	env := &core.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: ctxOrg.ID,
		Name:           req.Name,
		Slug:           req.Slug,
		Color:          req.Color,
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Environment().Create(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	env, _ = s.db.Environment().Get(ctx, database.EnvironmentByID(env.ID))

	return s.renderJSON(w, http.StatusOK, createEnvironmentResponse{
		Environment: environmentFromModel(env),
	})
}

type updateEnvironmentRequest struct {
	Name  *string `json:"name" validate:"required"`
	Color *string `json:"color" validate:"required"`
}

type updateEnvironmentResponse struct {
	Environment *environmentResponse `json:"environment"`
}

func (s *Server) handleUpdateEnvironment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	envIDReq := chi.URLParam(r, "environmentID")
	if envIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentID is required"))
	}

	var req updateEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditEnvironment); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	envID, err := uuid.FromString(envIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByOrganizationID(ctxOrg.ID), database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	if req.Name != nil {
		env.Name = internal.StringValue(req.Name)
	}
	if req.Color != nil {
		env.Color = internal.StringValue(req.Color)
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Environment().Update(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	env, _ = s.db.Environment().Get(ctx, database.EnvironmentByID(env.ID))

	return s.renderJSON(w, http.StatusOK, updateEnvironmentResponse{
		Environment: environmentFromModel(env),
	})
}

type deleteEnvironmentResponse struct {
	Environment *environmentResponse `json:"environment"`
}

func (s *Server) handleDeleteEnvironment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	envIDReq := chi.URLParam(r, "environmentID")
	if envIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentID is required"))
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditEnvironment); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	envID, err := uuid.FromString(envIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByOrganizationID(ctxOrg.ID), database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	if env.Slug == core.EnvironmentSlugDevelopment || env.Slug == core.EnvironmentSlugProduction {
		return errdefs.ErrInvalidArgument(errors.New("cannot delete development or production environment"))
	}

	apiKeys, err := s.db.APIKey().List(ctx, database.APIKeyByEnvironmentID(env.ID))
	if err != nil {
		return err
	}

	if len(apiKeys) > 0 {
		return errdefs.ErrEnvironmentDeletionNotAllowed(errors.New("cannot delete environment with API keys"))
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.Environment().Delete(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, deleteEnvironmentResponse{
		Environment: environmentFromModel(env),
	})
}
