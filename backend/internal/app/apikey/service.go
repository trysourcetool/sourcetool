package apikey

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
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
	currentOrg := ctxutil.CurrentOrganization(ctx)
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	apiKey, err := s.Repository.APIKey().Get(ctx, apikey.ByID(apiKeyID), apikey.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	return &dto.GetAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, env),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListAPIKeysOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	currentUser := ctxutil.CurrentUser(ctx)

	envs, err := s.Repository.Environment().List(ctx, environment.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	var devEnv *environment.Environment
	var liveEnvs []*environment.Environment
	for _, env := range envs {
		if env.Slug == environment.EnvironmentSlugDevelopment {
			devEnv = env
		} else {
			liveEnvs = append(liveEnvs, env)
		}
	}

	devKey, err := s.Repository.APIKey().Get(ctx, apikey.ByOrganizationID(currentOrg.ID), apikey.ByEnvironmentID(devEnv.ID), apikey.ByUserID(currentUser.ID))
	if err != nil {
		return nil, err
	}

	liveEnvIDs := make([]uuid.UUID, 0, len(liveEnvs))
	for _, env := range liveEnvs {
		liveEnvIDs = append(liveEnvIDs, env.ID)
	}
	liveKeys, err := s.Repository.APIKey().List(ctx, apikey.ByOrganizationID(currentOrg.ID), apikey.ByEnvironmentIDs(liveEnvIDs))
	if err != nil {
		return nil, err
	}

	liveKeyIDs := make([]uuid.UUID, 0, len(liveKeys))
	for _, key := range liveKeys {
		liveKeyIDs = append(liveKeyIDs, key.ID)
	}

	environments, err := s.Repository.Environment().MapByAPIKeyIDs(ctx, liveKeyIDs)
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
	currentOrg := ctxutil.CurrentOrganization(ctx)

	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	env, err := s.Repository.Environment().Get(ctx, environment.ByID(envID))
	if err != nil {
		return nil, err
	}

	if env.Slug == environment.EnvironmentSlugDevelopment {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot create API key for development environment"))
	}

	checker := permission.NewChecker(s.Repository)
	if env.Slug == environment.EnvironmentSlugDevelopment {
		if err := checker.AuthorizeOperation(ctx, permission.OperationEditDevModeAPIKey); err != nil {
			return nil, err
		}
	} else {
		if err := checker.AuthorizeOperation(ctx, permission.OperationEditLiveModeAPIKey); err != nil {
			return nil, err
		}
	}

	key, err := env.GenerateAPIKey()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	currentUser := ctxutil.CurrentUser(ctx)
	apiKey := &apikey.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: currentOrg.ID,
		EnvironmentID:  env.ID,
		UserID:         currentUser.ID,
		Name:           in.Name,
		Key:            key,
	}

	if err = s.Repository.RunTransaction(func(tx db.Transaction) error {
		return tx.APIKey().Create(ctx, apiKey)
	}); err != nil {
		return nil, err
	}

	apiKey, err = s.Repository.APIKey().Get(ctx, apikey.ByID(apiKey.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, nil),
	}, nil
}

func (s *ServiceCE) Update(ctx context.Context, in dto.UpdateAPIKeyInput) (*dto.UpdateAPIKeyOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	apiKeyID, err := uuid.FromString(in.APIKeyID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	apiKey, err := s.Repository.APIKey().Get(ctx, apikey.ByID(apiKeyID), apikey.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	if env.Slug == environment.EnvironmentSlugDevelopment {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot update API key for development environment"))
	}

	checker := permission.NewChecker(s.Repository)
	if env.Slug == environment.EnvironmentSlugDevelopment {
		if err := checker.AuthorizeOperation(ctx, permission.OperationEditDevModeAPIKey); err != nil {
			return nil, err
		}
	} else {
		if err := checker.AuthorizeOperation(ctx, permission.OperationEditLiveModeAPIKey); err != nil {
			return nil, err
		}
	}

	if in.Name != nil {
		apiKey.Name = ptrconv.SafeValue(in.Name)
	}

	if err = s.Repository.RunTransaction(func(tx db.Transaction) error {
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
	apiKey, err := s.Repository.APIKey().Get(ctx, apikey.ByID(apiKeyID))
	if err != nil {
		return nil, err
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByID(apiKey.EnvironmentID))
	if err != nil {
		return nil, err
	}

	if env.Slug == environment.EnvironmentSlugDevelopment {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot delete API key for development environment"))
	}

	checker := permission.NewChecker(s.Repository)
	if env.Slug == environment.EnvironmentSlugDevelopment {
		if err := checker.AuthorizeOperation(ctx, permission.OperationEditDevModeAPIKey); err != nil {
			return nil, err
		}
	} else {
		if err := checker.AuthorizeOperation(ctx, permission.OperationEditLiveModeAPIKey); err != nil {
			return nil, err
		}
	}

	if err = s.Repository.RunTransaction(func(tx db.Transaction) error {
		return tx.APIKey().Delete(ctx, apiKey)
	}); err != nil {
		return nil, err
	}

	return &dto.DeleteAPIKeyOutput{
		APIKey: dto.APIKeyFromModel(apiKey, nil),
	}, nil
}
