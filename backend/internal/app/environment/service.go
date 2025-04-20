package environment

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/app/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxdata"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	domainperm "github.com/trysourcetool/sourcetool/backend/internal/domain/permission"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

type Service interface {
	Get(context.Context, dto.GetEnvironmentInput) (*dto.GetEnvironmentOutput, error)
	List(context.Context) (*dto.ListEnvironmentsOutput, error)
	Create(context.Context, dto.CreateEnvironmentInput) (*dto.CreateEnvironmentOutput, error)
	Update(context.Context, dto.UpdateEnvironmentInput) (*dto.UpdateEnvironmentOutput, error)
	Delete(context.Context, dto.DeleteEnvironmentInput) (*dto.DeleteEnvironmentOutput, error)
}

type ServiceCE struct {
	*port.Dependencies
}

func NewServiceCE(d *port.Dependencies) *ServiceCE {
	return &ServiceCE{Dependencies: d}
}

func (s *ServiceCE) Get(ctx context.Context, in dto.GetEnvironmentInput) (*dto.GetEnvironmentOutput, error) {
	currentOrg := ctxdata.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByOrganizationID(currentOrg.ID), environment.ByID(envID))
	if err != nil {
		return nil, err
	}

	return &dto.GetEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListEnvironmentsOutput, error) {
	currentOrg := ctxdata.CurrentOrganization(ctx)
	envs, err := s.Repository.Environment().List(ctx, environment.ByOrganizationID(currentOrg.ID))
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
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, domainperm.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxdata.CurrentOrganization(ctx)

	slugExists, err := s.Repository.Environment().IsSlugExistsInOrganization(ctx, currentOrg.ID, in.Slug)
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, errdefs.ErrEnvironmentSlugAlreadyExists(errors.New("slug already exists"))
	}

	if !validateSlug(in.Slug) {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid slug"))
	}

	env := &environment.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: currentOrg.ID,
		Name:           in.Name,
		Slug:           in.Slug,
		Color:          in.Color,
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		if err := tx.Environment().Create(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	env, err = s.Repository.Environment().Get(ctx, environment.ByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}

func (s *ServiceCE) Update(ctx context.Context, in dto.UpdateEnvironmentInput) (*dto.UpdateEnvironmentOutput, error) {
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, domainperm.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxdata.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByID(envID), environment.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		env.Name = ptrconv.SafeValue(in.Name)
	}
	if in.Color != nil {
		env.Color = ptrconv.SafeValue(in.Color)
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		if err := tx.Environment().Update(ctx, env); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	env, err = s.Repository.Environment().Get(ctx, environment.ByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &dto.UpdateEnvironmentOutput{
		Environment: dto.EnvironmentFromModel(env),
	}, nil
}

func (s *ServiceCE) Delete(ctx context.Context, in dto.DeleteEnvironmentInput) (*dto.DeleteEnvironmentOutput, error) {
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, domainperm.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxdata.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByID(envID), environment.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if env.Slug == environment.EnvironmentSlugDevelopment || env.Slug == environment.EnvironmentSlugProduction {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot delete development or production environment"))
	}

	apiKeys, err := s.Repository.APIKey().List(ctx, apikey.ByEnvironmentID(env.ID))
	if err != nil {
		return nil, err
	}

	if len(apiKeys) > 0 {
		return nil, errdefs.ErrEnvironmentDeletionNotAllowed(errors.New("cannot delete environment with API keys"))
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
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
