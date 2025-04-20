package organization

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxdata"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

type Service interface {
	Create(ctx context.Context, in dto.CreateOrganizationInput) (*dto.CreateOrganizationOutput, error)
	CheckSubdomainAvailability(ctx context.Context, in dto.CheckSubdomainAvailabilityInput) error
}

type ServiceCE struct {
	*port.Dependencies
}

func NewServiceCE(d *port.Dependencies) *ServiceCE {
	return &ServiceCE{Dependencies: d}
}

func (s *ServiceCE) Create(ctx context.Context, in dto.CreateOrganizationInput) (*dto.CreateOrganizationOutput, error) {
	var subdomain *string
	if config.Config.IsCloudEdition {
		subdomain = ptrconv.NilValue(in.Subdomain)

		if lo.Contains(reservedSubdomains, in.Subdomain) {
			return nil, errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
		}

		exists, err := s.Repository.Organization().IsSubdomainExists(ctx, in.Subdomain)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
		}
	}

	o := &organization.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: subdomain,
	}

	currentUser := ctxdata.CurrentUser(ctx)

	orgAccess := &user.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         currentUser.ID,
		OrganizationID: o.ID,
		Role:           user.UserOrganizationRoleAdmin,
	}
	devEnv := &environment.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		Name:           environment.EnvironmentNameDevelopment,
		Slug:           environment.EnvironmentSlugDevelopment,
		Color:          environment.EnvironmentColorDevelopment,
	}
	envs := []*environment.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Name:           environment.EnvironmentNameProduction,
			Slug:           environment.EnvironmentSlugProduction,
			Color:          environment.EnvironmentColorProduction,
		},
		devEnv,
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}
	apiKey := &apikey.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: o.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         currentUser.ID,
		Name:           "",
		Key:            key,
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
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

	o, err = s.Repository.Organization().Get(ctx, organization.ByID(o.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateOrganizationOutput{
		Organization: dto.OrganizationFromModel(o),
	}, nil
}

func (s *ServiceCE) CheckSubdomainAvailability(ctx context.Context, in dto.CheckSubdomainAvailabilityInput) error {
	exists, err := s.Repository.Organization().IsSubdomainExists(ctx, in.Subdomain)
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
