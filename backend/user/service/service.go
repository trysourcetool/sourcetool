package service

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/jwt"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/permission"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	// User management methods
	GetMe(context.Context) (*output.GetMeOutput, error)
	UpdateMe(context.Context, input.UpdateMeInput) (*output.UpdateMeOutput, error)
	SendUpdateMeEmailInstructions(context.Context, input.SendUpdateMeEmailInstructionsInput) error
	UpdateMeEmail(context.Context, input.UpdateMeEmailInput) (*output.UpdateMeEmailOutput, error)

	// Organization methods
	List(context.Context) (*output.ListUsersOutput, error)
	Update(ctx context.Context, in input.UpdateUserInput) (*output.UpdateUserOutput, error)
	Delete(ctx context.Context, in input.DeleteUserInput) error

	// Invitation methods
	CreateUserInvitations(context.Context, input.CreateUserInvitationsInput) (*output.CreateUserInvitationsOutput, error)
	ResendUserInvitation(context.Context, input.ResendUserInvitationInput) (*output.ResendUserInvitationOutput, error)
}

// UserServiceCE implements the UserService interface for the Community Edition.
type UserServiceCE struct {
	*infra.Dependency
}

// NewUserServiceCE creates a new instance of the UserServiceCE.
func NewUserServiceCE(d *infra.Dependency) *UserServiceCE {
	return &UserServiceCE{Dependency: d}
}

func (s *UserServiceCE) GetMe(ctx context.Context) (*output.GetMeOutput, error) {
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByUserID(currentUser.ID),
		user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	role := orgAccess.Role

	return &output.GetMeOutput{
		User: output.UserFromModel(currentUser, currentOrg, role),
	}, nil
}

func (s *UserServiceCE) UpdateMe(ctx context.Context, in input.UpdateMeInput) (*output.UpdateMeOutput, error) {
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

	var role user.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &output.UpdateMeOutput{
		User: output.UserFromModel(currentUser, org, role),
	}, nil
}

// SendUpdateMeEmailInstructions sends instructions for updating a user's email address.
func (s *UserServiceCE) SendUpdateMeEmailInstructions(ctx context.Context, in input.SendUpdateMeEmailInstructionsInput) error {
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

	return s.Mailer.User().SendUpdateEmailInstructions(ctx, &user.SendUpdateUserEmailInstructions{
		To:        in.Email,
		FirstName: currentUser.FirstName,
		URL:       url,
	})
}

func (s *UserServiceCE) UpdateMeEmail(ctx context.Context, in input.UpdateMeEmailInput) (*output.UpdateMeEmailOutput, error) {
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
	u, err := s.Store.User().Get(ctx, user.ByID(userID))
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

	var role user.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &output.UpdateMeEmailOutput{
		User: output.UserFromModel(currentUser, org, role),
	}, nil
}

func (s *UserServiceCE) List(ctx context.Context) (*output.ListUsersOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)

	users, err := s.Store.User().List(ctx, user.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userInvitations, err := s.Store.User().ListInvitations(ctx, user.InvitationByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	roleMap := make(map[uuid.UUID]user.UserOrganizationRole)
	for _, oa := range orgAccesses {
		roleMap[oa.UserID] = oa.Role
	}

	usersOut := make([]*output.User, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, output.UserFromModel(u, nil, roleMap[u.ID]))
	}

	userInvitationsOut := make([]*output.UserInvitation, 0, len(userInvitations))
	for _, ui := range userInvitations {
		userInvitationsOut = append(userInvitationsOut, output.UserInvitationFromModel(ui))
	}

	return &output.ListUsersOutput{
		Users:           usersOut,
		UserInvitations: userInvitationsOut,
	}, nil
}

func (s *UserServiceCE) Update(ctx context.Context, in input.UpdateUserInput) (*output.UpdateUserOutput, error) {
	userID, err := uuid.FromString(in.UserID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	u, err := s.Store.User().Get(ctx, user.ByID(userID))
	if err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	if currentOrg == nil {
		return nil, errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, user.OrganizationAccessByOrganizationID(currentOrg.ID), user.OrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if in.Role != nil {
			orgAccess.Role = user.UserOrganizationRoleFromString(conv.SafeValue(in.Role))

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

	return &output.UpdateUserOutput{
		User: output.UserFromModel(u, currentOrg, orgAccess.Role),
	}, nil
}

func (s *UserServiceCE) Delete(ctx context.Context, in input.DeleteUserInput) error {
	checker := permission.NewChecker(s.Store)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditUser); err != nil {
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

	userToRemove, err := s.Store.User().Get(ctx, user.ByID(userIDToRemove))
	if err != nil {
		return err
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByUserID(userToRemove.ID),
		user.OrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		if errdefs.IsUserOrganizationAccessNotFound(err) {
			return nil
		}
		return err
	}

	if orgAccess.Role == user.UserOrganizationRoleAdmin {
		adminAccesses, err := s.Store.User().ListOrganizationAccesses(ctx,
			user.OrganizationAccessByOrganizationID(currentOrg.ID),
			user.OrganizationAccessByRole(user.UserOrganizationRoleAdmin))
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

func (s *UserServiceCE) CreateUserInvitations(ctx context.Context, in input.CreateUserInvitationsInput) (*output.CreateUserInvitationsOutput, error) {
	checker := permission.NewChecker(s.Store)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditUser); err != nil {
		return nil, err
	}

	o := ctxutil.CurrentOrganization(ctx)
	u := ctxutil.CurrentUser(ctx)

	invitations := make([]*user.UserInvitation, 0)
	emailInput := &user.SendInvitationEmail{
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

		invitations = append(invitations, &user.UserInvitation{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Email:          email,
			Role:           user.UserOrganizationRoleFromString(in.Role),
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

	usersInvitationsOut := make([]*output.UserInvitation, 0, len(invitations))
	for _, ui := range invitations {
		usersInvitationsOut = append(usersInvitationsOut, output.UserInvitationFromModel(ui))
	}

	return &output.CreateUserInvitationsOutput{
		UserInvitations: usersInvitationsOut,
	}, nil
}

func (s *UserServiceCE) ResendUserInvitation(ctx context.Context, in input.ResendUserInvitationInput) (*output.ResendUserInvitationOutput, error) {
	checker := permission.NewChecker(s.Store)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditUser); err != nil {
		return nil, err
	}

	invitationID, err := uuid.FromString(in.InvitationID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, user.InvitationByID(invitationID))
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

	emailInput := &user.SendInvitationEmail{
		Invitees: u.FullName(),
		URLs:     map[string]string{userInvitation.Email: url},
	}

	if err := s.Mailer.User().SendInvitationEmail(ctx, emailInput); err != nil {
		return nil, err
	}

	return &output.ResendUserInvitationOutput{
		UserInvitation: output.UserInvitationFromModel(userInvitation),
	}, nil
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *UserServiceCE) getUserOrganizationInfo(ctx context.Context) (*organization.Organization, *user.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, ctxutil.CurrentUser(ctx))
}

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *UserServiceCE) getOrganizationInfo(ctx context.Context, u *user.User) (*organization.Organization, *user.UserOrganizationAccess, error) {
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
func (s *UserServiceCE) getOrganizationBySubdomain(ctx context.Context, u *user.User, subdomain string) (*organization.Organization, *user.UserOrganizationAccess, error) {
	// Get organization by subdomain
	org, err := s.Store.Organization().Get(ctx, organization.BySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	// Verify user has access to this organization
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByOrganizationID(org.ID),
		user.OrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// getDefaultOrganizationForUser retrieves the default organization for a user
// (typically the most recently created one).
func (s *UserServiceCE) getDefaultOrganizationForUser(ctx context.Context, u *user.User) (*organization.Organization, *user.UserOrganizationAccess, error) {
	// Get user's organization access
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		user.OrganizationAccessByUserID(u.ID),
		user.OrganizationAccessOrderBy("created_at DESC"))
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.Store.Organization().Get(ctx, organization.ByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}
