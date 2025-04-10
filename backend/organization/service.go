package organization

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/authz"
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
	UpdateUser(ctx context.Context, in dto.UpdateOrganizationUserInput) (*dto.UpdateOrganizationUserOutput, error)
	DeleteUser(ctx context.Context, in dto.DeleteOrganizationUserInput) error
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

func (s *ServiceCE) UpdateUser(ctx context.Context, in dto.UpdateOrganizationUserInput) (*dto.UpdateOrganizationUserOutput, error) {
	userID, err := uuid.FromString(in.UserID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	u, err := s.Store.User().Get(ctx, storeopts.UserByID(userID))
	if err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg == nil {
		return nil, errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID), storeopts.UserOrganizationAccessByUserID(u.ID))
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

			existingGroups, err := tx.User().ListGroups(ctx, storeopts.UserGroupByUserID(u.ID))
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
		User: dto.UserFromModel(u, currentOrg, orgAccess.Role),
	}, nil
}

func (s *ServiceCE) DeleteUser(ctx context.Context, in dto.DeleteOrganizationUserInput) error {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditUser); err != nil {
		return err
	}

	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	userIDToRemove, err := uuid.FromString(in.UserID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if currentUser.ID == userIDToRemove {
		return errdefs.ErrPermissionDenied(errors.New("cannot remove yourself from the organization"))
	}

	userToRemove, err := s.Store.User().Get(ctx, storeopts.UserByID(userIDToRemove))
	if err != nil {
		return err
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByUserID(userToRemove.ID),
		storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		if errdefs.IsUserOrganizationAccessNotFound(err) {
			return nil
		}
		return err
	}

	if orgAccess.Role == model.UserOrganizationRoleAdmin {
		adminAccesses, err := s.Store.User().ListOrganizationAccesses(ctx,
			storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID),
			storeopts.UserOrganizationAccessByRole(int(model.UserOrganizationRoleAdmin)))
		if err != nil {
			return err
		}
		if len(adminAccesses) <= 1 {
			return errdefs.ErrPermissionDenied(errors.New("cannot remove the last admin from the organization"))
		}
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().DeleteOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		apiKeys, err := tx.APIKey().List(ctx, storeopts.APIKeyByUserID(userToRemove.ID), storeopts.APIKeyByOrganizationID(currentOrg.ID))
		if err != nil {
			return err
		}
		for _, apiKey := range apiKeys {
			if err := tx.APIKey().Delete(ctx, apiKey); err != nil {
				return err
			}
		}

		userGroups, err := tx.User().ListGroups(ctx, storeopts.UserGroupByUserID(userToRemove.ID), storeopts.UserGroupByOrganizationID(currentOrg.ID))
		if err != nil {
			return err
		}

		if len(userGroups) > 0 {
			if err := tx.User().BulkDeleteGroups(ctx, userGroups); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
