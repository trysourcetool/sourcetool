package user

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/app/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxdata"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	domainperm "github.com/trysourcetool/sourcetool/backend/internal/domain/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

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

type ServiceCE struct {
	*port.Dependencies
}

func NewServiceCE(d *port.Dependencies) *ServiceCE {
	return &ServiceCE{Dependencies: d}
}

func (s *ServiceCE) GetMe(ctx context.Context) (*dto.GetMeOutput, error) {
	currentUser := ctxdata.CurrentUser(ctx)
	currentOrg := ctxdata.CurrentOrganization(ctx)
	orgAccess, err := s.Repository.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByUserID(currentUser.ID),
		user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	role := orgAccess.Role

	return &dto.GetMeOutput{
		User: dto.UserFromModel(currentUser, currentOrg, role),
	}, nil
}

func (s *ServiceCE) UpdateMe(ctx context.Context, in dto.UpdateMeInput) (*dto.UpdateMeOutput, error) {
	currentUser := ctxdata.CurrentUser(ctx)

	if in.FirstName != nil {
		currentUser.FirstName = ptrconv.SafeValue(in.FirstName)
	}
	if in.LastName != nil {
		currentUser.LastName = ptrconv.SafeValue(in.LastName)
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		return tx.User().Update(ctx, currentUser)
	}); err != nil {
		return nil, err
	}

	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	var role user.UserOrganizationRole
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
	exists, err := s.Repository.User().IsEmailExists(ctx, in.Email)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(errors.New("email already exists"))
	}

	// Get current user and organization
	currentUser := ctxdata.CurrentUser(ctx)
	currentOrg := ctxdata.CurrentOrganization(ctx)

	// Create token for email update
	tok, err := createUpdateEmailToken(currentUser.ID.String(), in.Email)
	if err != nil {
		return err
	}

	// Build update URL
	url, err := buildUpdateEmailURL(ptrconv.SafeValue(currentOrg.Subdomain), tok)
	if err != nil {
		return err
	}

	return s.sendUpdateEmailInstructions(ctx, in.Email, currentUser.FirstName, url)
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
	u, err := s.Repository.User().Get(ctx, user.ByID(userID))
	if err != nil {
		return nil, err
	}

	currentUser := ctxdata.CurrentUser(ctx)
	if u.ID != currentUser.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("unauthorized"))
	}

	currentUser.Email = c.Email

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		return tx.User().Update(ctx, currentUser)
	}); err != nil {
		return nil, err
	}

	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	var role user.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.UpdateMeEmailOutput{
		User: dto.UserFromModel(currentUser, org, role),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListUsersOutput, error) {
	currentOrg := ctxdata.CurrentOrganization(ctx)

	users, err := s.Repository.User().List(ctx, user.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userInvitations, err := s.Repository.User().ListInvitations(ctx, user.InvitationByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	orgAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	roleMap := make(map[uuid.UUID]user.UserOrganizationRole)
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
	u, err := s.Repository.User().Get(ctx, user.ByID(userID))
	if err != nil {
		return nil, err
	}

	currentOrg := ctxdata.CurrentOrganization(ctx)
	if currentOrg == nil {
		return nil, errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	orgAccess, err := s.Repository.User().GetOrganizationAccess(ctx, user.OrganizationAccessByOrganizationID(currentOrg.ID), user.OrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		if in.Role != nil {
			orgAccess.Role = user.UserOrganizationRoleFromString(ptrconv.SafeValue(in.Role))

			if err := tx.User().UpdateOrganizationAccess(ctx, orgAccess); err != nil {
				return err
			}
		}

		if len(in.GroupIDs) != 0 {
			userGroups := make([]*user.UserGroup, 0, len(in.GroupIDs))
			for _, groupID := range in.GroupIDs {
				groupID, err := uuid.FromString(groupID)
				if err != nil {
					return err
				}
				userGroups = append(userGroups, &user.UserGroup{
					ID:      uuid.Must(uuid.NewV4()),
					UserID:  u.ID,
					GroupID: groupID,
				})
			}

			existingGroups, err := tx.User().ListGroups(ctx, user.GroupByUserID(u.ID))
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
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, domainperm.OperationEditUser); err != nil {
		return err
	}

	currentUser := ctxdata.CurrentUser(ctx)
	currentOrg := ctxdata.CurrentOrganization(ctx)
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

	userToRemove, err := s.Repository.User().Get(ctx, user.ByID(userIDToRemove))
	if err != nil {
		return err
	}

	orgAccess, err := s.Repository.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByUserID(userToRemove.ID),
		user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		if errdefs.IsUserOrganizationAccessNotFound(err) {
			return nil
		}
		return err
	}

	if orgAccess.Role == user.UserOrganizationRoleAdmin {
		adminAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx,
			user.OrganizationAccessByOrganizationID(currentOrg.ID),
			user.OrganizationAccessByRole(user.UserOrganizationRoleAdmin))
		if err != nil {
			return err
		}
		if len(adminAccesses) <= 1 {
			return errdefs.ErrPermissionDenied(errors.New("cannot remove the last admin from the organization"))
		}
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		if err := tx.User().DeleteOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		apiKeys, err := tx.APIKey().List(ctx, apikey.ByUserID(userToRemove.ID), apikey.ByOrganizationID(currentOrg.ID))
		if err != nil {
			return err
		}
		for _, apiKey := range apiKeys {
			if err := tx.APIKey().Delete(ctx, apiKey); err != nil {
				return err
			}
		}

		userGroups, err := tx.User().ListGroups(ctx, user.GroupByUserID(userToRemove.ID), user.GroupByOrganizationID(currentOrg.ID))
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
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, domainperm.OperationEditUser); err != nil {
		return nil, err
	}

	o := ctxdata.CurrentOrganization(ctx)
	u := ctxdata.CurrentUser(ctx)

	invitations := make([]*user.UserInvitation, 0)
	emailURLs := make(map[string]string)
	for _, email := range in.Emails {
		emailExsts, err := s.Repository.User().IsInvitationEmailExists(ctx, o.ID, email)
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

		url, err := buildInvitationURL(ptrconv.SafeValue(o.Subdomain), tok, email)
		if err != nil {
			return nil, err
		}

		emailURLs[email] = url

		invitations = append(invitations, &user.UserInvitation{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Email:          email,
			Role:           user.UserOrganizationRoleFromString(in.Role),
		})
	}

	if err := s.Repository.RunTransaction(func(tx port.Transaction) error {
		if err := tx.User().BulkInsertInvitations(ctx, invitations); err != nil {
			return err
		}

		if err := s.sendInvitationEmail(ctx, u.FullName(), emailURLs); err != nil {
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
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, domainperm.OperationEditUser); err != nil {
		return nil, err
	}

	invitationID, err := uuid.FromString(in.InvitationID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	userInvitation, err := s.Repository.User().GetInvitation(ctx, user.InvitationByID(invitationID))
	if err != nil {
		return nil, err
	}

	o := ctxdata.CurrentOrganization(ctx)
	if userInvitation.OrganizationID != o.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	u := ctxdata.CurrentUser(ctx)

	tok, err := createInvitationToken(userInvitation.Email)
	if err != nil {
		return nil, err
	}

	url, err := buildInvitationURL(ptrconv.SafeValue(o.Subdomain), tok, userInvitation.Email)
	if err != nil {
		return nil, err
	}

	emailURLs := map[string]string{userInvitation.Email: url}
	if err := s.sendInvitationEmail(ctx, u.FullName(), emailURLs); err != nil {
		return nil, err
	}

	return &dto.ResendUserInvitationOutput{
		UserInvitation: dto.UserInvitationFromModel(userInvitation),
	}, nil
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *ServiceCE) getUserOrganizationInfo(ctx context.Context) (*organization.Organization, *user.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, ctxdata.CurrentUser(ctx))
}

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *ServiceCE) getOrganizationInfo(ctx context.Context, u *user.User) (*organization.Organization, *user.UserOrganizationAccess, error) {
	if u == nil {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("user cannot be nil"))
	}

	subdomain := ctxdata.Subdomain(ctx)
	isCloudWithSubdomain := config.Config.IsCloudEdition && subdomain != "" && subdomain != "auth"

	// Different strategies for cloud vs. self-hosted or auth subdomain
	if isCloudWithSubdomain {
		return s.getOrganizationBySubdomain(ctx, u, subdomain)
	}

	return s.getDefaultOrganizationForUser(ctx, u)
}

// getOrganizationBySubdomain retrieves an organization by subdomain and verifies user access.
func (s *ServiceCE) getOrganizationBySubdomain(ctx context.Context, u *user.User, subdomain string) (*organization.Organization, *user.UserOrganizationAccess, error) {
	// Get organization by subdomain
	org, err := s.Repository.Organization().Get(ctx, organization.BySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	// Verify user has access to this organization
	orgAccess, err := s.Repository.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByOrganizationID(org.ID),
		user.OrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// getDefaultOrganizationForUser retrieves the default organization for a user
// (typically the most recently created one).
func (s *ServiceCE) getDefaultOrganizationForUser(ctx context.Context, u *user.User) (*organization.Organization, *user.UserOrganizationAccess, error) {
	// Get user's organization access
	orgAccess, err := s.Repository.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByUserID(u.ID),
		user.OrganizationAccessOrderBy("created_at DESC"))
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.Repository.Organization().Get(ctx, organization.ByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}
