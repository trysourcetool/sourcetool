package environment

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

type Service interface {
	Get(context.Context, types.GetEnvironmentInput) (*types.GetEnvironmentPayload, error)
	List(context.Context) (*types.ListEnvironmentsPayload, error)
	Create(context.Context, types.CreateEnvironmentInput) (*types.CreateEnvironmentPayload, error)
	Update(context.Context, types.UpdateEnvironmentInput) (*types.UpdateEnvironmentPayload, error)
	Delete(context.Context, types.DeleteEnvironmentInput) (*types.DeleteEnvironmentPayload, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Get(ctx context.Context, in types.GetEnvironmentInput) (*types.GetEnvironmentPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByOrganizationID(currentOrg.ID), model.EnvironmentByID(envID))
	if err != nil {
		return nil, err
	}

	return &types.GetEnvironmentPayload{
		Environment: &types.EnvironmentPayload{
			ID:        env.ID.String(),
			Name:      env.Name,
			Slug:      env.Slug,
			Color:     env.Color,
			CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*types.ListEnvironmentsPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	envs, err := s.Store.Environment().List(ctx, model.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	envsRes := make([]*types.EnvironmentPayload, 0, len(envs))
	for _, env := range envs {
		envsRes = append(envsRes, &types.EnvironmentPayload{
			ID:        env.ID.String(),
			Name:      env.Name,
			Slug:      env.Slug,
			Color:     env.Color,
			CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
		})
	}

	return &types.ListEnvironmentsPayload{
		Environments: envsRes,
	}, nil
}

func (s *ServiceCE) Create(ctx context.Context, in types.CreateEnvironmentInput) (*types.CreateEnvironmentPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)

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

	env, err = s.Store.Environment().Get(ctx, model.EnvironmentByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &types.CreateEnvironmentPayload{
		Environment: &types.EnvironmentPayload{
			ID:        env.ID.String(),
			Name:      env.Name,
			Slug:      env.Slug,
			Color:     env.Color,
			CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCE) Update(ctx context.Context, in types.UpdateEnvironmentInput) (*types.UpdateEnvironmentPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByID(envID), model.EnvironmentByOrganizationID(currentOrg.ID))
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

	env, err = s.Store.Environment().Get(ctx, model.EnvironmentByID(env.ID))
	if err != nil {
		return nil, err
	}

	return &types.UpdateEnvironmentPayload{
		Environment: &types.EnvironmentPayload{
			ID:        env.ID.String(),
			Name:      env.Name,
			Slug:      env.Slug,
			Color:     env.Color,
			CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCE) Delete(ctx context.Context, in types.DeleteEnvironmentInput) (*types.DeleteEnvironmentPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditEnvironment); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	env, err := s.Store.Environment().Get(ctx, model.EnvironmentByID(envID), model.EnvironmentByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if env.Slug == model.EnvironmentSlugDevelopment || env.Slug == model.EnvironmentSlugProduction {
		return nil, errdefs.ErrInvalidArgument(errors.New("cannot delete development or production environment"))
	}

	apiKeys, err := s.Store.APIKey().List(ctx, model.APIKeyByEnvironmentID(env.ID))
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

	return &types.DeleteEnvironmentPayload{
		Environment: &types.EnvironmentPayload{
			ID:        env.ID.String(),
			Name:      env.Name,
			Slug:      env.Slug,
			Color:     env.Color,
			CreatedAt: strconv.FormatInt(env.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(env.UpdatedAt.Unix(), 10),
		},
	}, nil
}
