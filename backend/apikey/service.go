package apikey

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type ServiceCE interface {
	Get(context.Context, types.GetAPIKeyInput) (*types.GetAPIKeyPayload, error)
	List(context.Context) (*types.ListAPIKeysPayload, error)
	Create(context.Context, types.CreateAPIKeyInput) (*types.CreateAPIKeyPayload, error)
	Update(context.Context, types.UpdateAPIKeyInput) (*types.UpdateAPIKeyPayload, error)
	Delete(context.Context, types.DeleteAPIKeyInput) (*types.DeleteAPIKeyPayload, error)
}

type ServiceCEImpl struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCEImpl {
	return &ServiceCEImpl{Dependency: d}
}

func (s *ServiceCEImpl) Get(ctx context.Context, in types.GetAPIKeyInput) (*types.GetAPIKeyPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	apiKey, err := s.Store.APIKey().Get(ctx, model.APIKeyByID(apiKeyID), model.APIKeyByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	return &types.GetAPIKeyPayload{
		APIKey: &types.APIKeyPayload{
			ID:        apiKey.ID.String(),
			Name:      apiKey.Name,
			Key:       apiKey.Key,
			CreatedAt: strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
			Environment: &types.EnvironmentPayload{
				ID:        env.ID.String(),
				Name:      env.Name,
				Slug:      env.Slug,
				Color:     env.Color,
				CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
				UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
			},
		},
	}, nil
}

func (s *ServiceCEImpl) List(ctx context.Context) (*types.ListAPIKeysPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	currentUser := ctxutils.CurrentUser(ctx)

	envs, err := s.Store.Environment().List(ctx, model.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	var devEnv *model.Environment
	var liveEnvs []*model.Environment
	for _, env := range envs {
		if env.Slug == model.EnvironmentSlugDevelopment {
			devEnv = env
		} else {
			liveEnvs = append(liveEnvs, env)
		}
	}

	devKey, err := s.Store.APIKey().Get(ctx, model.APIKeyByOrganizationID(currentOrg.ID), model.APIKeyByEnvironmentID(devEnv.ID), model.APIKeyByUserID(currentUser.ID))
	if err != nil {
		return nil, err
	}

	liveEnvIDs := make([]uuid.UUID, 0, len(liveEnvs))
	for _, env := range liveEnvs {
		liveEnvIDs = append(liveEnvIDs, env.ID)
	}
	liveKeys, err := s.Store.APIKey().List(ctx, model.APIKeyByOrganizationID(currentOrg.ID), model.APIKeyByEnvironmentIDs(liveEnvIDs))
	if err != nil {
		return nil, err
	}

	liveKeyIDs := make([]uuid.UUID, 0, len(liveKeys))
	for _, key := range liveKeys {
		liveKeyIDs = append(liveKeyIDs, key.ID)
	}

	environments, err := s.Store.Environment().MapByAPIKeyIDs(ctx, liveKeyIDs)
	if err != nil {
		return nil, err
	}

	liveKeysRes := make([]*types.APIKeyPayload, 0, len(liveKeys))
	for _, apiKey := range liveKeys {
		env, ok := environments[apiKey.ID]
		if !ok {
			return nil, errdefs.ErrEnvironmentNotFound(errors.New("environment not found"))
		}

		liveKeysRes = append(liveKeysRes, &types.APIKeyPayload{
			ID:        apiKey.ID.String(),
			Name:      apiKey.Name,
			Key:       apiKey.Key,
			CreatedAt: strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
			Environment: &types.EnvironmentPayload{
				ID:        env.ID.String(),
				Name:      env.Name,
				Slug:      env.Slug,
				Color:     env.Color,
				CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
				UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
			},
		})
	}

	return &types.ListAPIKeysPayload{
		DevKey: &types.APIKeyPayload{
			ID:        devKey.ID.String(),
			Name:      devKey.Name,
			Key:       devKey.Key,
			CreatedAt: strconv.FormatInt(devKey.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(devKey.UpdatedAt.Unix(), 10),
			Environment: &types.EnvironmentPayload{
				ID:        devEnv.ID.String(),
				Name:      devEnv.Name,
				Slug:      devEnv.Slug,
				Color:     devEnv.Color,
				CreatedAt: strconv.FormatInt(devEnv.CreatedAt.Unix(), 10),
				UpdatedAt: strconv.FormatInt(devEnv.UpdatedAt.Unix(), 10),
			},
		},
		LiveKeys: liveKeysRes,
	}, nil
}

func (s *ServiceCEImpl) Create(ctx context.Context, in types.CreateAPIKeyInput) (*types.CreateAPIKeyPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)

	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByID(envID))
	if err != nil {
		return nil, err
	}

	if env.Slug == model.EnvironmentSlugDevelopment {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot create API key for development environment"))
	}

	authorizer := authz.NewAuthorizer(s.Store)
	if env.Slug == model.EnvironmentSlugDevelopment {
		if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditDevModeAPIKey); err != nil {
			return nil, err
		}
	} else {
		if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditLiveModeAPIKey); err != nil {
			return nil, err
		}
	}

	key, err := model.GenerateAPIKey(currentOrg.Subdomain, env.Slug)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	currentUser := ctxutils.CurrentUser(ctx)
	apiKey := &model.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: currentOrg.ID,
		EnvironmentID:  env.ID,
		UserID:         currentUser.ID,
		Name:           in.Name,
		Key:            key,
	}

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.APIKey().Create(ctx, apiKey)
	}); err != nil {
		return nil, err
	}

	apiKey, err = s.Store.APIKey().Get(ctx, model.APIKeyByID(apiKey.ID))
	if err != nil {
		return nil, err
	}

	return &types.CreateAPIKeyPayload{
		APIKey: &types.APIKeyPayload{
			ID:        apiKey.ID.String(),
			Name:      apiKey.Name,
			Key:       apiKey.Key,
			CreatedAt: strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCEImpl) Update(ctx context.Context, in types.UpdateAPIKeyInput) (*types.UpdateAPIKeyPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	apiKey, err := s.Store.APIKey().Get(ctx, model.APIKeyByID(apiKeyID), model.APIKeyByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	if env.Slug == model.EnvironmentSlugDevelopment {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot update API key for development environment"))
	}

	authorizer := authz.NewAuthorizer(s.Store)
	if env.Slug == model.EnvironmentSlugDevelopment {
		if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditDevModeAPIKey); err != nil {
			return nil, err
		}
	} else {
		if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditLiveModeAPIKey); err != nil {
			return nil, err
		}
	}

	if in.Name != nil {
		apiKey.Name = conv.SafeValue(in.Name)
	}

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.APIKey().Update(ctx, apiKey)
	}); err != nil {
		return nil, err
	}

	return &types.UpdateAPIKeyPayload{
		APIKey: &types.APIKeyPayload{
			ID:        apiKey.ID.String(),
			Name:      apiKey.Name,
			Key:       apiKey.Key,
			CreatedAt: strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCEImpl) Delete(ctx context.Context, in types.DeleteAPIKeyInput) (*types.DeleteAPIKeyPayload, error) {
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	apiKey, err := s.Store.APIKey().Get(ctx, model.APIKeyByID(apiKeyID))
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	if env.Slug == model.EnvironmentSlugDevelopment {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot delete API key for development environment"))
	}

	authorizer := authz.NewAuthorizer(s.Store)
	if env.Slug == model.EnvironmentSlugDevelopment {
		if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditDevModeAPIKey); err != nil {
			return nil, err
		}
	} else {
		if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditLiveModeAPIKey); err != nil {
			return nil, err
		}
	}

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.APIKey().Delete(ctx, apiKey)
	}); err != nil {
		return nil, err
	}

	return &types.DeleteAPIKeyPayload{
		APIKey: &types.APIKeyPayload{
			ID:        apiKey.ID.String(),
			Name:      apiKey.Name,
			Key:       apiKey.Key,
			CreatedAt: strconv.FormatInt(apiKey.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(apiKey.UpdatedAt.Unix(), 10),
		},
	}, nil
}
