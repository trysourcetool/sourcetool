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

type apiKeyResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Key         string               `json:"key"`
	CreatedAt   string               `json:"createdAt"`
	UpdatedAt   string               `json:"updatedAt"`
	Environment *environmentResponse `json:"environment,omitempty"`
}

func (s *Server) apiKeyFromModel(apiKey *core.APIKey, env *core.Environment) *apiKeyResponse {
	if apiKey == nil {
		return nil
	}

	plainKey, err := s.encryptor.Decrypt(apiKey.KeyNonce, apiKey.KeyCiphertext)
	if err != nil {
		return nil
	}

	return &apiKeyResponse{
		ID:          apiKey.ID.String(),
		Name:        apiKey.Name,
		Key:         string(plainKey),
		CreatedAt:   strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
		UpdatedAt:   strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
		Environment: s.environmentFromModel(env),
	}
}

func (s *Server) hashAndEncryptAPIKey(plainAPIKey string) (keyHash string, keyNonce, keyCiphertext []byte, err error) {
	keyHash = core.HashAPIKey(plainAPIKey)
	keyNonce, keyCiphertext, err = s.encryptor.Encrypt([]byte(plainAPIKey))
	if err != nil {
		return "", nil, nil, err
	}

	return keyHash, keyNonce, keyCiphertext, nil
}

type getAPIKeyResponse struct {
	APIKey *apiKeyResponse `json:"apiKey"`
}

func (s *Server) handleGetAPIKey(w http.ResponseWriter, r *http.Request) error {
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

	return s.renderJSON(w, http.StatusOK, getAPIKeyResponse{
		APIKey: s.apiKeyFromModel(apiKey, env),
	})
}

type listAPIKeysResponse struct {
	DevKey   *apiKeyResponse   `json:"devKey"`
	LiveKeys []*apiKeyResponse `json:"liveKeys"`
}

func (s *Server) handleListAPIKeys(w http.ResponseWriter, r *http.Request) error {
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

	liveKeysOut := make([]*apiKeyResponse, 0, len(liveKeys))
	for _, apiKey := range liveKeys {
		env, ok := environments[apiKey.ID]
		if !ok {
			return errdefs.ErrEnvironmentNotFound(errors.New("environment not found"))
		}

		liveKeysOut = append(liveKeysOut, s.apiKeyFromModel(apiKey, env))
	}

	return s.renderJSON(w, http.StatusOK, listAPIKeysResponse{
		DevKey:   s.apiKeyFromModel(devKey, devEnv),
		LiveKeys: liveKeysOut,
	})
}

type createAPIKeyRequest struct {
	EnvironmentID string `json:"environmentId" validate:"required"`
	Name          string `json:"name" validate:"required"`
}

type createAPIKeyResponse struct {
	APIKey *apiKeyResponse `json:"apiKey"`
}

func (s *Server) handleCreateAPIKey(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req createAPIKeyRequest
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

	key, err := core.GenerateAPIKey(env.Slug)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	keyHash, keyNonce, keyCiphertext, err := s.hashAndEncryptAPIKey(key)
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
		KeyHash:        keyHash,
		KeyCiphertext:  keyCiphertext,
		KeyNonce:       keyNonce,
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

	return s.renderJSON(w, http.StatusOK, createAPIKeyResponse{
		APIKey: s.apiKeyFromModel(apiKey, env),
	})
}

type updateAPIKeyRequest struct {
	Name *string `json:"name" validate:"-"`
}

type updateAPIKeyResponse struct {
	APIKey *apiKeyResponse `json:"apiKey"`
}

func (s *Server) handleUpdateAPIKey(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	apiKeyIDReq := chi.URLParam(r, "apiKeyID")
	if apiKeyIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("apiKeyID is required"))
	}

	var req updateAPIKeyRequest
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

	return s.renderJSON(w, http.StatusOK, updateAPIKeyResponse{
		APIKey: s.apiKeyFromModel(apiKey, env),
	})
}

type deleteAPIKeyResponse struct {
	APIKey *apiKeyResponse `json:"apiKey"`
}

func (s *Server) handleDeleteAPIKey(w http.ResponseWriter, r *http.Request) error {
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

	return s.renderJSON(w, http.StatusOK, deleteAPIKeyResponse{
		APIKey: s.apiKeyFromModel(apiKey, env),
	})
}
