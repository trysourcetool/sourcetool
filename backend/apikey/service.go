package apikey

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type Service interface {
	Get(context.Context, dto.GetAPIKeyInput) (*dto.GetAPIKeyOutput, error)
	List(context.Context) (*dto.ListAPIKeysOutput, error)
	Create(context.Context, dto.CreateAPIKeyInput) (*dto.CreateAPIKeyOutput, error)
	Update(context.Context, dto.UpdateAPIKeyInput) (*dto.UpdateAPIKeyOutput, error)
	Delete(context.Context, dto.DeleteAPIKeyInput) (*dto.DeleteAPIKeyOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Get(ctx context.Context, in dto.GetAPIKeyInput) (*dto.GetAPIKeyOutput, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	apiKey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByID(apiKeyID), storeopts.APIKeyByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	return &dto.GetAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, env),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListAPIKeysOutput, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	currentUser := ctxutils.CurrentUser(ctx)

	envs, err := s.Store.Environment().List(ctx, storeopts.EnvironmentByOrganizationID(currentOrg.ID))
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

	devKey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByOrganizationID(currentOrg.ID), storeopts.APIKeyByEnvironmentID(devEnv.ID), storeopts.APIKeyByUserID(currentUser.ID))
	if err != nil {
		return nil, err
	}

	liveEnvIDs := make([]uuid.UUID, 0, len(liveEnvs))
	for _, env := range liveEnvs {
		liveEnvIDs = append(liveEnvIDs, env.ID)
	}
	liveKeys, err := s.Store.APIKey().List(ctx, storeopts.APIKeyByOrganizationID(currentOrg.ID), storeopts.APIKeyByEnvironmentIDs(liveEnvIDs))
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

	liveKeysOut := make([]*dto.APIKey, 0, len(liveKeys))
	for _, apiKey := range liveKeys {
		env, ok := environments[apiKey.ID]
		if !ok {
			return nil, errdefs.ErrEnvironmentNotFound(errors.New("environment not found"))
		}

		liveKeysOut = append(liveKeysOut, dto.APIKeyFromModel(apiKey, env))
	}

	return &dto.ListAPIKeysOutput{
		DevKey:   dto.APIKeyFromModel(devKey, devEnv),
		LiveKeys: liveKeysOut,
	}, nil
}

func (s *ServiceCE) Create(ctx context.Context, in dto.CreateAPIKeyInput) (*dto.CreateAPIKeyOutput, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)

	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(envID))
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

	key, err := env.GenerateAPIKey()
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

	apiKey, err = s.Store.APIKey().Get(ctx, storeopts.APIKeyByID(apiKey.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, nil),
	}, nil
}

func (s *ServiceCE) Update(ctx context.Context, in dto.UpdateAPIKeyInput) (*dto.UpdateAPIKeyOutput, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	apiKey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByID(apiKeyID), storeopts.APIKeyByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(apiKey.EnvironmentID))
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

	return &dto.UpdateAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, nil),
	}, nil
}

func (s *ServiceCE) Delete(ctx context.Context, in dto.DeleteAPIKeyInput) (*dto.DeleteAPIKeyOutput, error) {
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	apiKey, err := s.Store.APIKey().Get(ctx, storeopts.APIKeyByID(apiKeyID))
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(apiKey.EnvironmentID))
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

	return &dto.DeleteAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, nil),
	}, nil
}
