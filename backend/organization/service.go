package organization

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type Service interface {
	Create(ctx context.Context, in dto.CreateOrganizationInput) (*dto.CreateOrganizationOutput, error)
	CheckSubdomainAvailability(ctx context.Context, in dto.CheckSubdomainAvailabilityInput) error
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Create(ctx context.Context, in dto.CreateOrganizationInput) (*dto.CreateOrganizationOutput, error) {
	var subdomain *string
	if config.Config.IsCloudEdition {
		subdomain = conv.NilValue(in.Subdomain)

		if lo.Contains(reservedSubdomains, in.Subdomain) {
			return nil, errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
		}

		exists, err := s.Store.Organization().IsSubdomainExists(ctx, in.Subdomain)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
		}
	}

	o := &model.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: subdomain,
	}

	currentUser := ctxutil.CurrentUser(ctx)

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         currentUser.ID,
		OrganizationID: o.ID,
		Role:           model.UserOrganizationRoleAdmin,
	}
	devEnv := &model.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		Name:           model.EnvironmentNameDevelopment,
		Slug:           model.EnvironmentSlugDevelopment,
		Color:          model.EnvironmentColorDevelopment,
	}
	envs := []*model.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Name:           model.EnvironmentNameProduction,
			Slug:           model.EnvironmentSlugProduction,
			Color:          model.EnvironmentColorProduction,
		},
		devEnv,
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}
	apiKey := &model.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         currentUser.ID,
		Name:           "",
		Key:            key,
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Organization().Create(ctx, o); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		if err := tx.Environment().BulkInsert(ctx, envs); err != nil {
			return err
		}

		if err := tx.APIKey().Create(ctx, apiKey); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	o, err = s.Store.Organization().Get(ctx, storeopts.OrganizationByID(o.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateOrganizationOutput{
		Organization: dto.OrganizationFromModel(o),
	}, nil
}

func (s *ServiceCE) CheckSubdomainAvailability(ctx context.Context, in dto.CheckSubdomainAvailabilityInput) error {
	exists, err := s.Store.Organization().IsSubdomainExists(ctx, in.Subdomain)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
	}

	if lo.Contains(reservedSubdomains, in.Subdomain) {
		return errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
	}

	return nil
}
