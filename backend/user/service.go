package user

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/jwt"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

// Service defines the interface for user-related operations.
type Service interface {
	// User management methods
	GetMe(context.Context) (*dto.GetMeOutput, error)
	UpdateMe(context.Context, dto.UpdateMeInput) (*dto.UpdateMeOutput, error)
	SendUpdateMeEmailInstructions(context.Context, dto.SendUpdateMeEmailInstructionsInput) error
	UpdateMeEmail(context.Context, dto.UpdateMeEmailInput) (*dto.UpdateMeEmailOutput, error)

	// Organization methods
	List(context.Context) (*dto.ListUsersOutput, error)
	Update(ctx context.Context, in dto.UpdateUserInput) (*dto.UpdateUserOutput, error)
	Delete(ctx context.Context, in dto.DeleteUserInput) error

	// Invitation methods
	CreateUserInvitations(context.Context, dto.CreateUserInvitationsInput) (*dto.CreateUserInvitationsOutput, error)
	ResendUserInvitation(context.Context, dto.ResendUserInvitationInput) (*dto.ResendUserInvitationOutput, error)
}

// ServiceCE implements the Service interface for the Community Edition.
type ServiceCE struct {
	*infra.Dependency
}

// NewServiceCE creates a new instance of the ServiceCE.
func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) GetMe(ctx context.Context) (*dto.GetMeOutput, error) {
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByUserID(currentUser.ID),
		storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	role := orgAccess.Role

	return &dto.GetMeOutput{
		User: dto.UserFromModel(currentUser, currentOrg, role),
	}, nil
}

func (s *ServiceCE) UpdateMe(ctx context.Context, in dto.UpdateMeInput) (*dto.UpdateMeOutput, error) {
	currentUser := ctxutil.CurrentUser(ctx)

	if in.FirstName != nil {
		currentUser.FirstName = conv.SafeValue(in.FirstName)
	}
	if in.LastName != nil {
		currentUser.LastName = conv.SafeValue(in.LastName)
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, currentUser)
	}); err != nil {
		return nil, err
	}

	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	var role model.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.UpdateMeOutput{
		User: dto.UserFromModel(currentUser, org, role),
	}, nil
}

// SendUpdateMeEmailInstructions sends instructions for updating a user's email address.
func (s *ServiceCE) SendUpdateMeEmailInstructions(ctx context.Context, in dto.SendUpdateMeEmailInstructionsInput) error {
	// Validate email and confirmation match
	if in.Email != in.EmailConfirmation {
		return errdefs.ErrInvalidArgument(errors.New("email and email confirmation do not match"))
	}

	// Check if email already exists
	exists, err := s.Store.User().IsEmailExists(ctx, in.Email)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(errors.New("email already exists"))
	}

	// Get current user and organization
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)

	// Create token for email update
	tok, err := createUpdateEmailToken(currentUser.ID.String(), in.Email)
	if err != nil {
		return err
	}

	// Build update URL
	url, err := buildUpdateEmailURL(conv.SafeValue(currentOrg.Subdomain), tok)
	if err != nil {
		return err
	}

	return s.Mailer.User().SendUpdateEmailInstructions(ctx, &model.SendUpdateUserEmailInstructions{
		To:        in.Email,
		FirstName: currentUser.FirstName,
		URL:       url,
	})
}

func (s *ServiceCE) UpdateMeEmail(ctx context.Context, in dto.UpdateMeEmailInput) (*dto.UpdateMeEmailOutput, error) {
	c, err := jwt.ParseToken[*jwt.UserClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectUpdateEmail {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	u, err := s.Store.User().Get(ctx, storeopts.UserByID(userID))
	if err != nil {
		return nil, err
	}

	currentUser := ctxutil.CurrentUser(ctx)
	if u.ID != currentUser.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("unauthorized"))
	}

	currentUser.Email = c.Email

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, currentUser)
	}); err != nil {
		return nil, err
	}

	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	var role model.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.UpdateMeEmailOutput{
		User: dto.UserFromModel(currentUser, org, role),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListUsersOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)

	users, err := s.Store.User().List(ctx, storeopts.UserByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userInvitations, err := s.Store.User().ListInvitations(ctx, storeopts.UserInvitationByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	roleMap := make(map[uuid.UUID]model.UserOrganizationRole)
	for _, oa := range orgAccesses {
		roleMap[oa.UserID] = oa.Role
	}

	usersOut := make([]*dto.User, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, dto.UserFromModel(u, nil, roleMap[u.ID]))
	}

	userInvitationsOut := make([]*dto.UserInvitation, 0, len(userInvitations))
	for _, ui := range userInvitations {
		userInvitationsOut = append(userInvitationsOut, dto.UserInvitationFromModel(ui))
	}

	return &dto.ListUsersOutput{
		Users:           usersOut,
		UserInvitations: userInvitationsOut,
	}, nil
}

func (s *ServiceCE) Update(ctx context.Context, in dto.UpdateUserInput) (*dto.UpdateUserOutput, error) {
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

	return &dto.UpdateUserOutput{
		User: dto.UserFromModel(u, currentOrg, orgAccess.Role),
	}, nil
}

func (s *ServiceCE) Delete(ctx context.Context, in dto.DeleteUserInput) error {
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

func (s *ServiceCE) CreateUserInvitations(ctx context.Context, in dto.CreateUserInvitationsInput) (*dto.CreateUserInvitationsOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditUser); err != nil {
		return nil, err
	}

	o := ctxutil.CurrentOrganization(ctx)
	u := ctxutil.CurrentUser(ctx)

	invitations := make([]*model.UserInvitation, 0)
	emailInput := &model.SendInvitationEmail{
		Invitees: u.FullName(),
		URLs:     make(map[string]string),
	}
	for _, email := range in.Emails {
		emailExsts, err := s.Store.User().IsInvitationEmailExists(ctx, o.ID, email)
		if err != nil {
			return nil, err
		}
		if emailExsts {
			continue
		}

		tok, err := createInvitationToken(email)
		if err != nil {
			return nil, err
		}

		url, err := buildInvitationURL(conv.SafeValue(o.Subdomain), tok, email)
		if err != nil {
			return nil, err
		}

		emailInput.URLs[email] = url

		invitations = append(invitations, &model.UserInvitation{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Email:          email,
			Role:           model.UserOrganizationRoleFromString(in.Role),
		})
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().BulkInsertInvitations(ctx, invitations); err != nil {
			return err
		}

		if err := s.Mailer.User().SendInvitationEmail(ctx, emailInput); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	usersInvitationsOut := make([]*dto.UserInvitation, 0, len(invitations))
	for _, ui := range invitations {
		usersInvitationsOut = append(usersInvitationsOut, dto.UserInvitationFromModel(ui))
	}

	return &dto.CreateUserInvitationsOutput{
		UserInvitations: usersInvitationsOut,
	}, nil
}

func (s *ServiceCE) ResendUserInvitation(ctx context.Context, in dto.ResendUserInvitationInput) (*dto.ResendUserInvitationOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditUser); err != nil {
		return nil, err
	}

	invitationID, err := uuid.FromString(in.InvitationID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByID(invitationID))
	if err != nil {
		return nil, err
	}

	o := ctxutil.CurrentOrganization(ctx)
	if userInvitation.OrganizationID != o.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	u := ctxutil.CurrentUser(ctx)

	tok, err := createInvitationToken(userInvitation.Email)
	if err != nil {
		return nil, err
	}

	url, err := buildInvitationURL(conv.SafeValue(o.Subdomain), tok, userInvitation.Email)
	if err != nil {
		return nil, err
	}

	emailInput := &model.SendInvitationEmail{
		Invitees: u.FullName(),
		URLs:     map[string]string{userInvitation.Email: url},
	}

	if err := s.Mailer.User().SendInvitationEmail(ctx, emailInput); err != nil {
		return nil, err
	}

	return &dto.ResendUserInvitationOutput{
		UserInvitation: dto.UserInvitationFromModel(userInvitation),
	}, nil
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *ServiceCE) getUserOrganizationInfo(ctx context.Context) (*model.Organization, *model.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, ctxutil.CurrentUser(ctx))
}

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *ServiceCE) getOrganizationInfo(ctx context.Context, u *model.User) (*model.Organization, *model.UserOrganizationAccess, error) {
	if u == nil {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("user cannot be nil"))
	}

	subdomain := ctxutil.Subdomain(ctx)
	isCloudWithSubdomain := config.Config.IsCloudEdition && subdomain != "" && subdomain != "auth"

	// Different strategies for cloud vs. self-hosted or auth subdomain
	if isCloudWithSubdomain {
		return s.getOrganizationBySubdomain(ctx, u, subdomain)
	}

	return s.getDefaultOrganizationForUser(ctx, u)
}

// getOrganizationBySubdomain retrieves an organization by subdomain and verifies user access.
func (s *ServiceCE) getOrganizationBySubdomain(ctx context.Context, u *model.User, subdomain string) (*model.Organization, *model.UserOrganizationAccess, error) {
	// Get organization by subdomain
	org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	// Verify user has access to this organization
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByOrganizationID(org.ID),
		storeopts.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// (typically the most recently created one).
func (s *ServiceCE) getDefaultOrganizationForUser(ctx context.Context, u *model.User) (*model.Organization, *model.UserOrganizationAccess, error) {
	// Get user's organization access
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByUserID(u.ID),
		storeopts.UserOrganizationAccessOrderBy("created_at DESC"))
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}
