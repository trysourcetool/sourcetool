package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) getAPIKey(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	apiKeyIDReq := chi.URLParam(r, "apiKeyID")
	if apiKeyIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("apiKeyID is required"))
	}

	apiKeyID, err := uuid.FromString(apiKeyIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(apiKeyID))
	if err != nil {
		return err
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.GetAPIKeyResponse{
		APIKey: responses.APIKeyFromModel(apiKey, env),
	})
}

func (s *Server) listAPIKeys(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	currentOrg := internal.ContextOrganization(ctx)
	ctxUser := internal.ContextUser(ctx)

	envs, err := s.db.Environment().List(ctx, database.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	var devEnv *core.Environment
	var liveEnvs []*core.Environment
	for _, env := range envs {
		if env.Slug == core.EnvironmentSlugDevelopment {
			devEnv = env
		} else {
			liveEnvs = append(liveEnvs, env)
		}
	}

	devKey, err := s.db.APIKey().Get(ctx, database.APIKeyByOrganizationID(currentOrg.ID), database.APIKeyByEnvironmentID(devEnv.ID), database.APIKeyByUserID(ctxUser.ID))
	if err != nil {
		return err
	}

	liveEnvIDs := make([]uuid.UUID, 0, len(liveEnvs))
	for _, env := range liveEnvs {
		liveEnvIDs = append(liveEnvIDs, env.ID)
	}
	liveKeys, err := s.db.APIKey().List(ctx, database.APIKeyByOrganizationID(currentOrg.ID), database.APIKeyByEnvironmentIDs(liveEnvIDs))
	if err != nil {
		return err
	}

	liveKeyIDs := make([]uuid.UUID, 0, len(liveKeys))
	for _, key := range liveKeys {
		liveKeyIDs = append(liveKeyIDs, key.ID)
	}

	environments, err := s.db.Environment().MapByAPIKeyIDs(ctx, liveKeyIDs)
	if err != nil {
		return err
	}

	liveKeysOut := make([]*responses.APIKeyResponse, 0, len(liveKeys))
	for _, apiKey := range liveKeys {
		env, ok := environments[apiKey.ID]
		if !ok {
			return errdefs.ErrEnvironmentNotFound(errors.New("environment not found"))
		}

		liveKeysOut = append(liveKeysOut, responses.APIKeyFromModel(apiKey, env))
	}

	return s.renderJSON(w, http.StatusOK, responses.ListAPIKeysResponse{
		DevKey:   responses.APIKeyFromModel(devKey, devEnv),
		LiveKeys: liveKeysOut,
	})
}

func (s *Server) createAPIKey(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)

	envID, err := uuid.FromString(req.EnvironmentID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	if env.Slug == core.EnvironmentSlugDevelopment {
		return errdefs.ErrInvalidArgument(errors.New("cannot create API key for development environment"))
	}

	if env.Slug == core.EnvironmentSlugDevelopment {
		if err := s.checker.AuthorizeOperation(ctx, core.OperationEditDevModeAPIKey); err != nil {
			return err
		}
	} else {
		if err := s.checker.AuthorizeOperation(ctx, core.OperationEditLiveModeAPIKey); err != nil {
			return err
		}
	}

	apiKeys, err := s.db.APIKey().List(ctx, database.APIKeyByOrganizationID(ctxOrg.ID), database.APIKeyByEnvironmentID(env.ID))
	if err != nil {
		return err
	}

	if len(apiKeys) >= 1 {
		// currently, we only support one API key per environment.
		// we will support multiple API keys per environment in the future:
		// concern:
		// 1. since multiple host instances can be associated with a single session, a client must call initializeClient multiple times; otherwise, there will be hosts without a session.
		// 2. when closing a session, all related host instances must be properly closed as well.
		// 3. additionally, if the same session ID is reused, the session persistence logic in initializeClient will also need to be updated.
		return errdefs.ErrInvalidArgument(errors.New("cannot create more than one API key for this environment"))
	}

	key, err := env.GenerateAPIKey()
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	ctxUser := internal.ContextUser(ctx)
	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: ctxOrg.ID,
		EnvironmentID:  env.ID,
		UserID:         ctxUser.ID,
		Name:           req.Name,
		Key:            key,
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.APIKey().Create(ctx, apiKey); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	apiKey, _ = s.db.APIKey().Get(ctx, database.APIKeyByID(apiKey.ID))

	return s.renderJSON(w, http.StatusOK, responses.CreateAPIKeyResponse{
		APIKey: responses.APIKeyFromModel(apiKey, env),
	})
}

func (s *Server) updateAPIKey(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	apiKeyIDReq := chi.URLParam(r, "apiKeyID")
	if apiKeyIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("apiKeyID is required"))
	}

	var req requests.UpdateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	apiKeyID, err := uuid.FromString(apiKeyIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(apiKeyID), database.APIKeyByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return err
	}

	if env.Slug == core.EnvironmentSlugDevelopment {
		return errdefs.ErrInvalidArgument(errors.New("cannot update API key for development environment"))
	}

	if env.Slug == core.EnvironmentSlugDevelopment {
		if err := s.checker.AuthorizeOperation(ctx, core.OperationEditDevModeAPIKey); err != nil {
			return err
		}
	} else {
		if err := s.checker.AuthorizeOperation(ctx, core.OperationEditLiveModeAPIKey); err != nil {
			return err
		}
	}

	if req.Name != nil {
		apiKey.Name = internal.StringValue(req.Name)
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.APIKey().Update(ctx, apiKey); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.UpdateAPIKeyResponse{
		APIKey: responses.APIKeyFromModel(apiKey, env),
	})
}

func (s *Server) deleteAPIKey(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	apiKeyIDReq := chi.URLParam(r, "apiKeyID")
	if apiKeyIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("apiKeyID is required"))
	}

	apiKeyID, err := uuid.FromString(apiKeyIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	apiKey, err := s.db.APIKey().Get(ctx, database.APIKeyByID(apiKeyID))
	if err != nil {
		return err
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return err
	}

	if env.Slug == core.EnvironmentSlugDevelopment {
		return errdefs.ErrInvalidArgument(errors.New("cannot delete API key for development environment"))
	}

	if env.Slug == core.EnvironmentSlugDevelopment {
		if err := s.checker.AuthorizeOperation(ctx, core.OperationEditDevModeAPIKey); err != nil {
			return err
		}
	} else {
		if err := s.checker.AuthorizeOperation(ctx, core.OperationEditLiveModeAPIKey); err != nil {
			return err
		}
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.APIKey().Delete(ctx, apiKey); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.DeleteAPIKeyResponse{
		APIKey: responses.APIKeyFromModel(apiKey, env),
	})
}
