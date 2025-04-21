package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/internal/mail"
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

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

func createUpdateEmailToken(userID, email string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(core.EmailTokenExpiration)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectUpdateEmail,
		},
	})
}

func createInvitationToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(core.EmailTokenExpiration)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectInvitation,
		},
	})
}

func (s *Server) sendUpdateEmailInstructions(ctx context.Context, to, firstName, url string) error {
	subject := "[Sourcetool] Confirm your new email address"
	content := fmt.Sprintf(`Hi %s,

We received a request to change the email address associated with your Sourcetool account. To ensure the security of your account, we need you to verify your new email address.

Please click the following link within the next 24 hours to confirm your email change:
%s

Thank you for being a part of the Sourcetool community!
Regards,

Sourcetool Team`,
		firstName,
		url,
	)

	return mail.Send(ctx, mail.MailInput{
		From:     config.Config.SMTP.FromEmail,
		FromName: "Sourcetool Team",
		To:       []string{to},
		Subject:  subject,
		Body:     content,
	})
}

func (s *Server) sendInvitationEmail(ctx context.Context, invitees string, emaiURLs map[string]string) error {
	subject := "You've been invited to join Sourcetool!"
	baseContent := `You've been invited to join Sourcetool!

To accept the invitation, please create your account by clicking the URL below within 24 hours.

%s

- This URL will expire in 24 hours.
- This is a send-only email address.
- Your account will not be created unless you complete the next steps.`

	sendEmails := make([]string, 0)
	for email, url := range emaiURLs {
		if lo.Contains(sendEmails, email) {
			continue
		}

		content := fmt.Sprintf(baseContent, url)

		if err := mail.Send(ctx, mail.MailInput{
			From:     config.Config.SMTP.FromEmail,
			FromName: "Sourcetool Team",
			To:       []string{email},
			Subject:  subject,
			Body:     content,
		}); err != nil {
			return err
		}

		sendEmails = append(sendEmails, email)
	}

	return nil
}

func (s *Server) getMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	currentUser := internal.CurrentUser(ctx)
	currentOrg := internal.CurrentOrganization(ctx)
	orgAccess, err := s.db.User().GetOrganizationAccess(ctx,
		database.UserOrganizationAccessByUserID(currentUser.ID),
		database.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}
	role := orgAccess.Role

	return s.renderJSON(w, http.StatusOK, responses.GetMeResponse{
		User: responses.UserFromModel(currentUser, role, currentOrg),
	})
}

func (s *Server) updateMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	currentUser := internal.CurrentUser(ctx)

	if req.FirstName != nil {
		currentUser.FirstName = internal.SafeValue(req.FirstName)
	}
	if req.LastName != nil {
		currentUser.LastName = internal.SafeValue(req.LastName)
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, currentUser); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	org, orgAccess, err := s.resolveOrganization(ctx, currentUser)
	if err != nil {
		return err
	}

	var role core.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return s.renderJSON(w, http.StatusOK, responses.UpdateMeResponse{
		User: responses.UserFromModel(currentUser, role, org),
	})
}

func (s *Server) sendUpdateMeEmailInstructions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.SendUpdateMeEmailInstructionsRequest
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
	currentUser := internal.CurrentUser(ctx)
	currentOrg := internal.CurrentOrganization(ctx)

	// Create token for email update
	tok, err := createUpdateEmailToken(currentUser.ID.String(), req.Email)
	if err != nil {
		return err
	}

	// Build update URL
	url, err := buildUpdateEmailURL(internal.SafeValue(currentOrg.Subdomain), tok)
	if err != nil {
		return err
	}

	return s.sendUpdateEmailInstructions(ctx, req.Email, currentUser.FirstName, url)
}

func (s *Server) updateMeEmail(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.UpdateMeEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	c, err := jwt.ParseToken[*jwt.UserClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if c.Subject != jwt.UserSignatureSubjectUpdateEmail {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	u, err := s.db.User().Get(ctx, database.UserByID(userID))
	if err != nil {
		return err
	}

	currentUser := internal.CurrentUser(ctx)
	if u.ID != currentUser.ID {
		return errdefs.ErrUnauthenticated(errors.New("unauthorized"))
	}

	currentUser.Email = c.Email

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, currentUser); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	org, orgAccess, err := s.resolveOrganization(ctx, currentUser)
	if err != nil {
		return err
	}

	var role core.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return s.renderJSON(w, http.StatusOK, responses.UpdateMeEmailResponse{
		User: responses.UserFromModel(currentUser, role, org),
	})
}

func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	currentOrg := internal.CurrentOrganization(ctx)

	users, err := s.db.User().List(ctx, database.UserByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	userInvitations, err := s.db.User().ListInvitations(ctx, database.UserInvitationByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}

	orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return err
	}
	roleMap := make(map[uuid.UUID]core.UserOrganizationRole)
	for _, oa := range orgAccesses {
		roleMap[oa.UserID] = oa.Role
	}

	usersOut := make([]*responses.UserResponse, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, responses.UserFromModel(u, roleMap[u.ID], currentOrg))
	}

	userInvitationsOut := make([]*responses.UserInvitationResponse, 0, len(userInvitations))
	for _, ui := range userInvitations {
		userInvitationsOut = append(userInvitationsOut, responses.UserInvitationFromModel(ui))
	}

	return s.renderJSON(w, http.StatusOK, responses.ListUsersResponse{
		Users:           usersOut,
		UserInvitations: userInvitationsOut,
	})
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userIDReq := chi.URLParam(r, "userID")
	if userIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("userID is required"))
	}

	userID, err := uuid.FromString(userIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	var req requests.UpdateUserRequest
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

	currentOrg := internal.CurrentOrganization(ctx)
	if currentOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	orgAccess, err := s.db.User().GetOrganizationAccess(ctx,
		database.UserOrganizationAccessByOrganizationID(currentOrg.ID),
		database.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return err
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if req.Role != nil {
			orgAccess.Role = core.UserOrganizationRoleFromString(internal.SafeValue(req.Role))

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

	return s.renderJSON(w, http.StatusOK, responses.UpdateUserResponse{
		User: responses.UserFromModel(u, orgAccess.Role, currentOrg),
	})
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userIDReq := chi.URLParam(r, "userID")
	if userIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("userID is required"))
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	currentUser := internal.CurrentUser(ctx)
	currentOrg := internal.CurrentOrganization(ctx)
	if currentOrg == nil {
		return errdefs.ErrUnauthenticated(errors.New("current organization not found"))
	}

	userIDToRemove, err := uuid.FromString(userIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if currentUser.ID == userIDToRemove {
		return errdefs.ErrPermissionDenied(errors.New("cannot remove yourself from the organization"))
	}

	userToRemove, err := s.db.User().Get(ctx, database.UserByID(userIDToRemove))
	if err != nil {
		return err
	}

	orgAccess, err := s.db.User().GetOrganizationAccess(ctx,
		database.UserOrganizationAccessByUserID(userToRemove.ID),
		database.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		if errdefs.IsUserOrganizationAccessNotFound(err) {
			return nil
		}
		return err
	}

	if orgAccess.Role == core.UserOrganizationRoleAdmin {
		adminAccesses, err := s.db.User().ListOrganizationAccesses(ctx,
			database.UserOrganizationAccessByOrganizationID(currentOrg.ID),
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

		apiKeys, err := s.db.APIKey().List(ctx, database.APIKeyByUserID(userToRemove.ID), database.APIKeyByOrganizationID(currentOrg.ID))
		if err != nil {
			return err
		}
		for _, apiKey := range apiKeys {
			if err := tx.APIKey().Delete(ctx, apiKey); err != nil {
				return err
			}
		}

		userGroups, err := s.db.User().ListGroups(ctx, database.UserGroupByUserID(userToRemove.ID), database.UserGroupByOrganizationID(currentOrg.ID))
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

	return s.renderJSON(w, http.StatusOK, responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully deleted user",
	})
}

func (s *Server) createUserInvitations(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.CreateUserInvitationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	o := internal.CurrentOrganization(ctx)
	u := internal.CurrentUser(ctx)

	invitations := make([]*core.UserInvitation, 0)
	emailURLs := make(map[string]string)
	for _, email := range req.Emails {
		emailExsts, err := s.db.User().IsInvitationEmailExists(ctx, o.ID, email)
		if err != nil {
			return err
		}
		if emailExsts {
			continue
		}

		tok, err := createInvitationToken(email)
		if err != nil {
			return err
		}

		url, err := buildInvitationURL(internal.SafeValue(o.Subdomain), tok, email)
		if err != nil {
			return err
		}

		emailURLs[email] = url

		invitations = append(invitations, &core.UserInvitation{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Email:          email,
			Role:           core.UserOrganizationRoleFromString(req.Role),
		})
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().BulkInsertInvitations(ctx, invitations); err != nil {
			return err
		}

		if err := s.sendInvitationEmail(ctx, u.FullName(), emailURLs); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	usersInvitationsOut := make([]*responses.UserInvitationResponse, 0, len(invitations))
	for _, ui := range invitations {
		usersInvitationsOut = append(usersInvitationsOut, responses.UserInvitationFromModel(ui))
	}

	return s.renderJSON(w, http.StatusOK, responses.CreateUserInvitationsResponse{
		UserInvitations: usersInvitationsOut,
	})
}

func (s *Server) resendUserInvitation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	invitationIDReq := chi.URLParam(r, "invitationID")
	if invitationIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("invitationID is required"))
	}

	invitationID, err := uuid.FromString(invitationIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := s.checker.AuthorizeOperation(ctx, core.OperationEditUser); err != nil {
		return err
	}

	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByID(invitationID))
	if err != nil {
		return err
	}

	o := internal.CurrentOrganization(ctx)
	if userInvitation.OrganizationID != o.ID {
		return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	u := internal.CurrentUser(ctx)

	tok, err := createInvitationToken(userInvitation.Email)
	if err != nil {
		return err
	}

	url, err := buildInvitationURL(internal.SafeValue(o.Subdomain), tok, userInvitation.Email)
	if err != nil {
		return err
	}

	emailURLs := map[string]string{userInvitation.Email: url}
	if err := s.sendInvitationEmail(ctx, u.FullName(), emailURLs); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.ResendUserInvitationResponse{
		UserInvitation: responses.UserInvitationFromModel(userInvitation),
	})
}
