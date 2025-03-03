package organization

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type Service interface {
	Create(ctx context.Context, in dto.CreateOrganizationInput) (*dto.CreateOrganizationOutput, error)
	CheckSubdomainAvailability(ctx context.Context, in dto.CheckSubdomainAvailabilityInput) error
	UpdateUser(ctx context.Context, in dto.UpdateOrganizationUserInput) (*dto.UpdateOrganizationUserOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) Create(ctx context.Context, in dto.CreateOrganizationInput) (*dto.CreateOrganizationOutput, error) {
	exists, err := s.Store.Organization().IsSubdomainExists(ctx, in.Subdomain)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain already exists"))
	}

	if lo.Contains(reservedSubdomains, in.Subdomain) {
		return nil, errdefs.ErrOrganizationSubdomainAlreadyExists(errors.New("subdomain is reserved"))
	}

	o := &model.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: in.Subdomain,
	}

	currentUser := ctxutils.CurrentUser(ctx)

	// TODO: currentUserがすでに組織に所属しているかをチェック

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         currentUser.ID,
		OrganizationID: o.ID,
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
	key, err := model.GenerateAPIKey(o.Subdomain, devEnv.Slug)
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

	o, err = s.Store.Organization().Get(ctx, model.OrganizationByID(o.ID))
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

func (s *ServiceCE) UpdateUser(ctx context.Context, in dto.UpdateOrganizationUserInput) (*dto.UpdateOrganizationUserOutput, error) {
	userID, err := uuid.FromString(in.UserID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	u, err := s.Store.User().Get(ctx, model.UserByID(userID))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	o, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByOrganizationID(o.ID), model.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if in.Role != nil {
			orgAccess.Role = model.UserOrganizationRoleFromString(conv.SafeValue(in.Role))

			if err := tx.User().UpdateOrganizationAccess(ctx, orgAccess); err != nil {
				return err
			}
		}

		if len(in.GroupIDs) != 0 {
			userGroups := make([]*model.UserGroup, 0, len(in.GroupIDs))
			for _, groupID := range in.GroupIDs {
				groupID, err := uuid.FromString(groupID)
				if err != nil {
					return err
				}
				userGroups = append(userGroups, &model.UserGroup{
					ID:      uuid.Must(uuid.NewV4()),
					UserID:  u.ID,
					GroupID: groupID,
				})
			}

			existingGroups, err := tx.User().ListGroups(ctx, model.UserGroupByUserID(u.ID))
			if err != nil {
				return err
			}

			if err := tx.User().BulkDeleteGroups(ctx, existingGroups); err != nil {
				return err
			}

			if err := tx.User().BulkInsertGroups(ctx, userGroups); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.UpdateOrganizationUserOutput{
		User: dto.UserFromModel(u, o, orgAccess.Role),
	}, nil
}
