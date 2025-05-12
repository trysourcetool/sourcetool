package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/internal/mail"
)

type userResponse struct {
	ID           string                `json:"id"`
	Email        string                `json:"email"`
	FirstName    string                `json:"firstName"`
	LastName     string                `json:"lastName"`
	Role         string                `json:"role"`
	CreatedAt    string                `json:"createdAt"`
	UpdatedAt    string                `json:"updatedAt"`
	Organization *organizationResponse `json:"organization"`
}

func (s *Server) userFromModel(user *core.User, role core.UserOrganizationRole, org *core.Organization) *userResponse {
	if user == nil {
		return nil
	}

	return &userResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         role.String(),
		CreatedAt:    strconv.FormatInt(user.CreatedAt.Unix(), 10),
		UpdatedAt:    strconv.FormatInt(user.UpdatedAt.Unix(), 10),
		Organization: s.organizationFromModel(org),
	}
}

type userInvitationResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

func (s *Server) userInvitationFromModel(invitation *core.UserInvitation) *userInvitationResponse {
	return &userInvitationResponse{
		ID:        invitation.ID.String(),
		Email:     invitation.Email,
		CreatedAt: strconv.FormatInt(invitation.CreatedAt.Unix(), 10),
	}
}

type userGroupResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	GroupID   string `json:"groupId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (s *Server) userGroupFromModel(userGroup *core.UserGroup) *userGroupResponse {
	return &userGroupResponse{
		ID:        userGroup.ID.String(),
		UserID:    userGroup.UserID.String(),
		GroupID:   userGroup.GroupID.String(),
		CreatedAt: strconv.FormatInt(userGroup.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(userGroup.UpdatedAt.Unix(), 10),
	}
}

func buildUpdateEmailURL(subdomain, token string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "email", "update", "confirm"), map[string]string{
		"token": token,
	})
}

func buildInvitationURL(subdomain, token, email string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("auth", "invitations", "login"), map[string]string{
		"token": token,
		"email": email,
	})
}

type getMeResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	ctxUser := internal.ContextUser(ctx)
	ctxOrg := internal.ContextOrganization(ctx)
	orgAccess, err := s.db.User().GetOrganizationAccess(ctx,
		database.UserOrganizationAccessByUserID(ctxUser.ID),
		database.UserOrganizationAccessByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}
	role := orgAccess.Role

	return s.renderJSON(w, http.StatusOK, getMeResponse{
		User: s.userFromModel(ctxUser, role, ctxOrg),
	})
}

type updateMeRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type updateMeResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleUpdateMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req updateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	ctxUser := internal.ContextUser(ctx)

	if req.FirstName != nil {
		ctxUser.FirstName = internal.StringValue(req.FirstName)
	}
	if req.LastName != nil {
		ctxUser.LastName = internal.StringValue(req.LastName)
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, ctxUser); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	org, orgAccess, err := s.resolveOrganization(ctx, ctxUser)
	if err != nil {
		return err
	}

	var role core.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return s.renderJSON(w, http.StatusOK, updateMeResponse{
		User: s.userFromModel(ctxUser, role, org),
	})
}

type sendUpdateMeEmailInstructionsRequest struct {
	Email             string `json:"email" validate:"required,email"`
	EmailConfirmation string `json:"emailConfirmation" validate:"required,email"`
}

func (s *Server) handleSendUpdateMeEmailInstructions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req sendUpdateMeEmailInstructionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Validate email and confirmation match
	if req.Email != req.EmailConfirmation {
		return errdefs.ErrInvalidArgument(errors.New("email and email confirmation do not match"))
	}

	// Check if email already exists
	exists, err := s.db.User().IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(errors.New("email already exists"))
	}

	// Get current user and organization
	ctxUser := internal.ContextUser(ctx)
	ctxOrg := internal.ContextOrganization(ctx)

	// Create token for email update
	tok, err := jwt.SignUpdateUserEmailToken(ctxUser.ID.String(), req.Email)
	if err != nil {
		return err
	}

	// Build update URL
	url, err := buildUpdateEmailURL(internal.StringValue(ctxOrg.Subdomain), tok)
	if err != nil {
		return err
	}

	if err := mail.SendUpdateEmailInstructions(ctx, req.Email, ctxUser.FirstName, url); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, statusResponse{
		Code:    http.StatusOK,
		Message: "Email update instructions sent successfully",
	})
}

type updateMeEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type updateMeEmailResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleUpdateMeEmail(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req updateMeEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	c, err := jwt.ParseUpdateUserEmailClaims(req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	userID, err := uuid.FromString(c.Subject)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	u, err := s.db.User().Get(ctx, database.UserByID(userID))
	if err != nil {
		return err
	}

	ctxUser := internal.ContextUser(ctx)
	if u.ID != ctxUser.ID {
		return errdefs.ErrUnauthenticated(errors.New("unauthorized"))
	}

	ctxUser.Email = c.Email

	if ctxUser.GoogleID != "" {
		ctxUser.GoogleID = ""
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, ctxUser); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	org, orgAccess, err := s.resolveOrganization(ctx, ctxUser)
	if err != nil {
		return err
	}

	var role core.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return s.renderJSON(w, http.StatusOK, updateMeEmailResponse{
		User: s.userFromModel(ctxUser, role, org),
	})
}

type listUsersResponse struct {
	Users           []*userResponse           `json:"users"`
	UserInvitations []*userInvitationResponse `json:"userInvitations"`
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	ctxOrg := internal.ContextOrganization(ctx)

	users, err := s.db.User().List(ctx, database.UserByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	userInvitations, err := s.db.User().ListInvitations(ctx, database.UserInvitationByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}

	orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByOrganizationID(ctxOrg.ID))
	if err != nil {
		return err
	}
	roleMap := make(map[uuid.UUID]core.UserOrganizationRole)
	for _, oa := range orgAccesses {
		roleMap[oa.UserID] = oa.Role
	}

	usersOut := make([]*userResponse, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, s.userFromModel(u, roleMap[u.ID], ctxOrg))
	}

	userInvitationsOut := make([]*userInvitationResponse, 0, len(userInvitations))
	for _, ui := range userInvitations {
		userInvitationsOut = append(userInvitationsOut, s.userInvitationFromModel(ui))
	}

	return s.renderJSON(w, http.StatusOK, listUsersResponse{
		Users:           usersOut,
		UserInvitations: userInvitationsOut,
	})
}

type updateUserRequest struct {
	Role     *string  `json:"role" validate:"oneof=admin developer member"`
	GroupIDs []string `json:"groupIds"`
}

type updateUserResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userIDReq := chi.URLParam(r, "userID")
	if userIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("userID is required"))
	}

	userID, err := uuid.FromString(userIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	u, err := s.db.User().Get(ctx, database.UserByID(userID))
	if err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	if ctxOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	orgAccess, err := s.db.User().GetOrganizationAccess(ctx,
		database.UserOrganizationAccessByOrganizationID(ctxOrg.ID),
		database.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return err
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if req.Role != nil {
			orgAccess.Role = core.UserOrganizationRoleFromString(internal.StringValue(req.Role))

			if err := tx.User().UpdateOrganizationAccess(ctx, orgAccess); err != nil {
				return err
			}
		}

		if len(req.GroupIDs) != 0 {
			userGroups := make([]*core.UserGroup, 0, len(req.GroupIDs))
			for _, groupID := range req.GroupIDs {
				groupID, err := uuid.FromString(groupID)
				if err != nil {
					return err
				}
				userGroups = append(userGroups, &core.UserGroup{
					ID:      uuid.Must(uuid.NewV4()),
					UserID:  u.ID,
					GroupID: groupID,
				})
			}

			existingGroups, err := s.db.User().ListGroups(ctx, database.UserGroupByUserID(u.ID))
			if err != nil {
				return err
			}

			if err := s.db.User().BulkDeleteGroups(ctx, existingGroups); err != nil {
				return err
			}

			if err := s.db.User().BulkInsertGroups(ctx, userGroups); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, updateUserResponse{
		User: s.userFromModel(u, orgAccess.Role, ctxOrg),
	})
}

type deleteUserResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userIDReq := chi.URLParam(r, "userID")
	if userIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("userID is required"))
	}

	if err := s.permissionChecker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	ctxUser := internal.ContextUser(ctx)
	ctxOrg := internal.ContextOrganization(ctx)
	if ctxOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	userIDToRemove, err := uuid.FromString(userIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if ctxUser.ID == userIDToRemove {
		return errdefs.ErrPermissionDenied(errors.New("cannot remove yourself from the organization"))
	}

	userToRemove, err := s.db.User().Get(ctx, database.UserByID(userIDToRemove))
	if err != nil {
		return err
	}

	orgAccess, err := s.db.User().GetOrganizationAccess(ctx,
		database.UserOrganizationAccessByUserID(userToRemove.ID),
		database.UserOrganizationAccessByOrganizationID(ctxOrg.ID))
	if err != nil {
		if errdefs.IsUserOrganizationAccessNotFound(err) {
			return nil
		}
		return err
	}

	if orgAccess.Role == core.UserOrganizationRoleAdmin {
		adminAccesses, err := s.db.User().ListOrganizationAccesses(ctx,
			database.UserOrganizationAccessByOrganizationID(ctxOrg.ID),
			database.UserOrganizationAccessByRole(core.UserOrganizationRoleAdmin))
		if err != nil {
			return err
		}
		if len(adminAccesses) <= 1 {
			return errdefs.ErrPermissionDenied(errors.New("cannot remove the last admin from the organization"))
		}
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().DeleteOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		apiKeys, err := s.db.APIKey().List(ctx, database.APIKeyByUserID(userToRemove.ID), database.APIKeyByOrganizationID(ctxOrg.ID))
		if err != nil {
			return err
		}
		for _, apiKey := range apiKeys {
			if err := tx.APIKey().Delete(ctx, apiKey); err != nil {
				return err
			}
		}

		userGroups, err := s.db.User().ListGroups(ctx, database.UserGroupByUserID(userToRemove.ID), database.UserGroupByOrganizationID(ctxOrg.ID))
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

	return s.renderJSON(w, http.StatusOK, deleteUserResponse{
		User: s.userFromModel(userToRemove, orgAccess.Role, ctxOrg),
	})
}

type createUserInvitationsRequest struct {
	Emails []string `json:"emails" validate:"required"`
	Role   string   `json:"role" validate:"required,oneof=admin developer member"`
}

type createUserInvitationsResponse struct {
	UserInvitations []*userInvitationResponse `json:"userInvitations"`
}

func (s *Server) handleCreateUserInvitations(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req createUserInvitationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.permissionChecker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	ctxUser := internal.ContextUser(ctx)

	validEmails := make([]string, 0, len(req.Emails))
	for _, email := range req.Emails {
		emailExists, err := s.db.User().IsInvitationEmailExists(ctx, ctxOrg.ID, email)
		if err != nil {
			return err
		}
		if !emailExists {
			validEmails = append(validEmails, email)
		}
	}

	// Check if we can add all these users to the organization (CE limit check)
	if err := s.canAddUsersToOrganization(ctx, ctxOrg.ID, len(validEmails)); err != nil {
		return err
	}

	invitations := make([]*core.UserInvitation, 0, len(validEmails))
	emailURLs := make(map[string]string)
	for _, email := range validEmails {
		tok, err := jwt.SignInvitationToken(email)
		if err != nil {
			return err
		}

		url, err := buildInvitationURL(internal.StringValue(ctxOrg.Subdomain), tok, email)
		if err != nil {
			return err
		}

		emailURLs[email] = url

		invitations = append(invitations, &core.UserInvitation{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: ctxOrg.ID,
			Email:          email,
			Role:           core.UserOrganizationRoleFromString(req.Role),
		})
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().BulkInsertInvitations(ctx, invitations); err != nil {
			return err
		}

		if err := mail.SendInvitationEmail(ctx, ctxUser.FullName(), emailURLs); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	usersInvitationsOut := make([]*userInvitationResponse, 0, len(invitations))
	for _, ui := range invitations {
		usersInvitationsOut = append(usersInvitationsOut, s.userInvitationFromModel(ui))
	}

	return s.renderJSON(w, http.StatusOK, createUserInvitationsResponse{
		UserInvitations: usersInvitationsOut,
	})
}

type resendUserInvitationResponse struct {
	UserInvitation *userInvitationResponse `json:"userInvitation"`
}

func (s *Server) handleResendUserInvitation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	invitationIDReq := chi.URLParam(r, "invitationID")
	if invitationIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("invitationID is required"))
	}

	invitationID, err := uuid.FromString(invitationIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := s.permissionChecker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByID(invitationID))
	if err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	if userInvitation.OrganizationID != ctxOrg.ID {
		return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	ctxUser := internal.ContextUser(ctx)

	tok, err := jwt.SignInvitationToken(userInvitation.Email)
	if err != nil {
		return err
	}

	url, err := buildInvitationURL(internal.StringValue(ctxOrg.Subdomain), tok, userInvitation.Email)
	if err != nil {
		return err
	}

	emailURLs := map[string]string{userInvitation.Email: url}
	if err := mail.SendInvitationEmail(ctx, ctxUser.FullName(), emailURLs); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, resendUserInvitationResponse{
		UserInvitation: s.userInvitationFromModel(userInvitation),
	})
}

type deleteUserInvitationResponse struct {
	UserInvitation *userInvitationResponse `json:"userInvitation"`
}

func (s *Server) handleDeleteUserInvitation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	invitationIDReq := chi.URLParam(r, "invitationID")
	if invitationIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("invitationID is required"))
	}

	invitationID, err := uuid.FromString(invitationIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := s.permissionChecker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByID(invitationID))
	if err != nil {
		return err
	}

	ctxOrg := internal.ContextOrganization(ctx)
	if userInvitation.OrganizationID != ctxOrg.ID {
		return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	if err := s.db.User().DeleteInvitation(ctx, userInvitation); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, deleteUserInvitationResponse{
		UserInvitation: s.userInvitationFromModel(userInvitation),
	})
}
