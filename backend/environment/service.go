package environment

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type Service interface {
	Get(context.Context, dto.GetEnvironmentInput) (*dto.GetEnvironmentOutput, error)
	List(context.Context) (*dto.ListEnvironmentsOutput, error)
	Create(context.Context, dto.CreateEnvironmentInput) (*dto.CreateEnvironmentOutput, error)
	Update(context.Context, dto.UpdateEnvironmentInput) (*dto.UpdateEnvironmentOutput, error)
	Delete(context.Context, dto.DeleteEnvironmentInput) (*dto.DeleteEnvironmentOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Get(ctx context.Context, in dto.GetEnvironmentInput) (*dto.GetEnvironmentOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByOrganizationID(currentOrg.ID), storeopts.EnvironmentByID(envID))
	if err != nil {
		return nil, err
	}

	return &dto.GetEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListEnvironmentsOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	envs, err := s.Store.Environment().List(ctx, storeopts.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	envsOut := make([]*dto.Environment, 0, len(envs))
	for _, env := range envs {
		envsOut = append(envsOut, dto.EnvironmentFromModel(env))
	}

	return &dto.ListEnvironmentsOutput{
		Environments: envsOut,
	}, nil
}

func (s *ServiceCE) Create(ctx context.Context, in dto.CreateEnvironmentInput) (*dto.CreateEnvironmentOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)

	slugExists, err := s.Store.Environment().IsSlugExistsInOrganization(ctx, currentOrg.ID, in.Slug)
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, errdefs.ErrEnvironmentSlugAlreadyExists(errors.New("slug already exists"))
	}

	if !validateSlug(in.Slug) {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid slug"))
	}

	env := &model.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: currentOrg.ID,
		Name:           in.Name,
		Slug:           in.Slug,
		Color:          in.Color,
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Environment().Create(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	env, err = s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}

func (s *ServiceCE) Update(ctx context.Context, in dto.UpdateEnvironmentInput) (*dto.UpdateEnvironmentOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(envID), storeopts.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		env.Name = conv.SafeValue(in.Name)
	}
	if in.Color != nil {
		env.Color = conv.SafeValue(in.Color)
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Environment().Update(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	env, err = s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &dto.UpdateEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}

func (s *ServiceCE) Delete(ctx context.Context, in dto.DeleteEnvironmentInput) (*dto.DeleteEnvironmentOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(envID), storeopts.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if env.Slug == model.EnvironmentSlugDevelopment || env.Slug == model.EnvironmentSlugProduction {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot delete development or production environment"))
	}

	apiKeys, err := s.Store.APIKey().List(ctx, storeopts.APIKeyByEnvironmentID(env.ID))
	if err != nil {
		return nil, err
	}

	if len(apiKeys) > 0 {
		return nil, errdefs.ErrEnvironmentDeletionNotAllowed(errors.New("cannot delete environment with API keys"))
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Environment().Delete(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.DeleteEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}
