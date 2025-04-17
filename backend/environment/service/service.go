package service

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/permission"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type EnvironmentService interface {
	Get(context.Context, input.GetEnvironmentInput) (*output.GetEnvironmentOutput, error)
	List(context.Context) (*output.ListEnvironmentsOutput, error)
	Create(context.Context, input.CreateEnvironmentInput) (*output.CreateEnvironmentOutput, error)
	Update(context.Context, input.UpdateEnvironmentInput) (*output.UpdateEnvironmentOutput, error)
	Delete(context.Context, input.DeleteEnvironmentInput) (*output.DeleteEnvironmentOutput, error)
}

type EnvironmentServiceCE struct {
	*infra.Dependency
}

func NewEnvironmentServiceCE(d *infra.Dependency) *EnvironmentServiceCE {
	return &EnvironmentServiceCE{Dependency: d}
}

func (s *EnvironmentServiceCE) Get(ctx context.Context, in input.GetEnvironmentInput) (*output.GetEnvironmentOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, environment.ByOrganizationID(currentOrg.ID), environment.ByID(envID))
	if err != nil {
		return nil, err
	}

	return &output.GetEnvironmentOutput{
		Environment: output.EnvironmentFromModel(env),
	}, nil
}

func (s *EnvironmentServiceCE) List(ctx context.Context) (*output.ListEnvironmentsOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	envs, err := s.Store.Environment().List(ctx, environment.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	envsOut := make([]*output.Environment, 0, len(envs))
	for _, env := range envs {
		envsOut = append(envsOut, output.EnvironmentFromModel(env))
	}

	return &output.ListEnvironmentsOutput{
		Environments: envsOut,
	}, nil
}

func (s *EnvironmentServiceCE) Create(ctx context.Context, in input.CreateEnvironmentInput) (*output.CreateEnvironmentOutput, error) {
	checker := permission.NewChecker(s.Store)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditEnvironment); err != nil {
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

	env := &environment.Environment{
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

	env, err = s.Store.Environment().Get(ctx, environment.ByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &output.CreateEnvironmentOutput{
		Environment: output.EnvironmentFromModel(env),
	}, nil
}

func (s *EnvironmentServiceCE) Update(ctx context.Context, in input.UpdateEnvironmentInput) (*output.UpdateEnvironmentOutput, error) {
	checker := permission.NewChecker(s.Store)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, environment.ByID(envID), environment.ByOrganizationID(currentOrg.ID))
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

	env, err = s.Store.Environment().Get(ctx, environment.ByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &output.UpdateEnvironmentOutput{
		Environment: output.EnvironmentFromModel(env),
	}, nil
}

func (s *EnvironmentServiceCE) Delete(ctx context.Context, in input.DeleteEnvironmentInput) (*output.DeleteEnvironmentOutput, error) {
	checker := permission.NewChecker(s.Store)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, environment.ByID(envID), environment.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if env.Slug == environment.EnvironmentSlugDevelopment || env.Slug == environment.EnvironmentSlugProduction {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot delete development or production environment"))
	}

	apiKeys, err := s.Store.APIKey().List(ctx, apikey.ByEnvironmentID(env.ID))
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

	return &output.DeleteEnvironmentOutput{
		Environment: output.EnvironmentFromModel(env),
	}, nil
}
