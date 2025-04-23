package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
	gojwt "github.com/golang-jwt/jwt/v5"

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

func createMagicLinkToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserEmailClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectMagicLink,
		},
	})
}

func createInvitationMagicLinkToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserEmailClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectInvitationMagicLink,
		},
	})
}

func createMagicLinkRegistrationToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserEmailClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectMagicLinkRegistration,
		},
	})
}

func buildMagicLinkURL(subdomain, token string) (string, error) {
	base := config.Config.AuthBaseURL()
	if subdomain != "" && subdomain != "auth" {
		base = config.Config.OrgBaseURL(subdomain)
	}
	return internal.BuildURL(base, path.Join("auth", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}

func buildInvitationMagicLinkURL(subdomain, token string) (string, error) {
	baseURL := config.Config.OrgBaseURL(subdomain)
	return internal.BuildURL(baseURL, path.Join("auth", "invitations", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}

func (s *Server) sendMagicLinkEmail(ctx context.Context, email, firstName, url string) error {
	subject := "Log in to your Sourcetool account"

	content := fmt.Sprintf(`Hi %s,

Here's your magic link to log in to your Sourcetool account. Click the link below to access your account securely without a password:

%s

- This link will expire in 15 minutes for security reasons.
- If you didn't request this link, you can safely ignore this email.

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, url)

	if err := mail.Send(ctx, mail.MailInput{
		From:     config.Config.SMTP.FromEmail,
		FromName: "Sourcetool Team",
		To:       []string{email},
		Subject:  subject,
		Body:     content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Server) sendInvitationMagicLinkEmail(ctx context.Context, email, firstName, url string) error {
	subject := "Your invitation to join Sourcetool"

	content := fmt.Sprintf(`Hi %s,

You've been invited to join Sourcetool. Click the link below to accept the invitation:

%s

This link will expire in 15 minutes.

Best regards,
The Sourcetool Team`, firstName, url)

	if err := mail.Send(ctx, mail.MailInput{
		From:     config.Config.SMTP.FromEmail,
		FromName: "Sourcetool Team",
		To:       []string{email},
		Subject:  subject,
		Body:     content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *Server) sendMultipleOrganizationsMagicLinkEmail(ctx context.Context, email, firstName string, loginURLs []string) error {
	subject := "Choose your Sourcetool organization to log in"

	urlList := ""
	for _, url := range loginURLs {
		urlList += url + "\n"
	}

	content := fmt.Sprintf(`Hi %s,

Your email, %s, is associated with multiple Sourcetool organizations. You may log in to each one by clicking its magic link below:

%s

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, email, urlList)

	if err := mail.Send(ctx, mail.MailInput{
		From:     config.Config.SMTP.FromEmail,
		FromName: "Sourcetool Team",
		To:       []string{email},
		Subject:  subject,
		Body:     content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Server) requestMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.RequestMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Check if email exists
	exists, err := s.db.User().IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}

	var firstName string
	isNewUser := !exists

	// Handle Cloud Edition with subdomain
	if config.Config.IsCloudEdition {
		ctxSubdomain := internal.ContextSubdomain(ctx)
		if ctxSubdomain != "" && ctxSubdomain != "auth" {
			// Get organization by subdomain
			org, err := s.db.Organization().Get(ctx, database.OrganizationBySubdomain(ctxSubdomain))
			if err != nil {
				return err
			}

			if exists {
				// For existing users, check if they have access to this organization
				u, err := s.db.User().Get(ctx, database.UserByEmail(req.Email))
				if err != nil {
					return err
				}

				_, err = s.db.User().GetOrganizationAccess(ctx,
					database.UserOrganizationAccessByUserID(u.ID),
					database.UserOrganizationAccessByOrganizationID(org.ID))
				if err != nil {
					return errdefs.ErrUnauthenticated(errors.New("user does not have access to this organization"))
				}
			} else {
				// For new users, registration is only allowed through invitations
				return errdefs.ErrPermissionDenied(errors.New("registration is only allowed through invitations"))
			}
		}
	}

	if exists {
		// Get user by email for existing users
		u, err := s.db.User().Get(ctx, database.UserByEmail(req.Email))
		if err != nil {
			return err
		}
		firstName = u.FirstName

		// Get user's organization access information
		orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByUserID(u.ID))
		if err != nil {
			return err
		}

		// Cloud edition specific handling for multiple organizations
		if config.Config.IsCloudEdition && len(orgAccesses) > 1 {
			// Handle multiple organizations
			loginURLs := make([]string, 0, len(orgAccesses))
			for _, access := range orgAccesses {
				org, err := s.db.Organization().Get(ctx, database.OrganizationByID(access.OrganizationID))
				if err != nil {
					return err
				}

				// Create org-specific magic link
				tok, err := createMagicLinkToken(req.Email)
				if err != nil {
					return err
				}

				url, err := buildMagicLinkURL(internal.SafeValue(org.Subdomain), tok)
				if err != nil {
					return err
				}
				loginURLs = append(loginURLs, url)
			}

			if err := s.sendMultipleOrganizationsMagicLinkEmail(ctx, req.Email, firstName, loginURLs); err != nil {
				return err
			}

			return s.renderJSON(w, http.StatusOK, &responses.RequestMagicLinkResponse{
				Email: req.Email,
				IsNew: false,
			})
		}
	} else {
		// For new users, generate a temporary ID that will be verified/used later
		firstName = "there" // Default greeting

		// For self-hosted mode, check if creating an organization is allowed
		if !config.Config.IsCloudEdition {
			// Check if an organization already exists in self-hosted mode
			if err := s.validateSelfHostedOrganization(ctx); err != nil {
				return err
			}
		}
	}

	// Determine subdomain context based on edition
	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = internal.ContextSubdomain(ctx)
	}

	// Create token for magic link authentication
	tok, err := createMagicLinkToken(req.Email)
	if err != nil {
		return err
	}

	// Build magic link URL
	url, err := buildMagicLinkURL(subdomain, tok)
	if err != nil {
		return err
	}

	// Send magic link email
	if err := s.sendMagicLinkEmail(ctx, req.Email, firstName, url); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, responses.RequestMagicLinkResponse{
		Email: req.Email,
		IsNew: isNewUser,
	})
}

func (s *Server) authenticateWithMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.AuthenticateWithMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	c, err := jwt.ParseToken[*jwt.UserEmailClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if c.Subject != jwt.UserSignatureSubjectMagicLink {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Check if user exists
	exists, err := s.db.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := createMagicLinkRegistrationToken(c.Email)
		if err != nil {
			return fmt.Errorf("failed to generate registration token: %w", err)
		}

		return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithMagicLinkResponse{
			Token:           registrationToken,
			IsNewUser:       true,
			HasOrganization: false,
		})
	}

	// Get existing user
	u, err := s.db.User().Get(ctx, database.UserByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get user's organization access information
	orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return err
	}

	// Handle organization subdomain logic
	ctxSubdomain := internal.ContextSubdomain(ctx)
	var orgAccess *core.UserOrganizationAccess
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if ctxSubdomain != "auth" {
			// For specific organization subdomain, resolve org and access
			orgAccess, err = s.db.User().GetOrganizationAccess(ctx,
				database.UserOrganizationAccessByUserID(u.ID),
				database.UserOrganizationAccessByOrganizationSubdomain(ctxSubdomain))
			if err != nil {
				return err
			}
			orgSubdomain = ctxSubdomain
		} else {
			// For auth subdomain
			if len(orgAccesses) == 0 {
				// No organization - sign in as a user not associated with any organization
			} else if len(orgAccesses) == 1 {
				// Single organization - redirect to it
				orgAccess = orgAccesses[0]
				org, err := s.db.Organization().Get(ctx, database.OrganizationByID(orgAccess.OrganizationID))
				if err != nil {
					return err
				}
				orgSubdomain = internal.SafeValue(org.Subdomain)
			} else {
				return errdefs.ErrUserMultipleOrganizations(errors.New("user has multiple organizations"))
			}
		}
	} else {
		// Self-hosted mode has only one organization
		orgAccess = orgAccesses[0]
		_, err = s.db.Organization().Get(ctx, database.OrganizationByID(orgAccess.OrganizationID))
		if err != nil {
			return err
		}
	}

	// Create token, refresh token, etc.
	token, xsrfToken, _, hashedRefreshToken, _, err := s.createTokens(
		u.ID, core.TmpTokenExpiration)
	if err != nil {
		return err
	}

	// Update user with new refresh token
	u.RefreshTokenHash = hashedRefreshToken
	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return err
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, u); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	hasOrganization := orgAccess != nil
	if !hasOrganization {
		cookieConfig := newCookieConfig()
		cookieConfig.SetTmpAuthCookie(w, token, xsrfToken, config.Config.AuthDomain())
	}

	return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithMagicLinkResponse{
		AuthURL:         authURL,
		Token:           token,
		HasOrganization: hasOrganization,
		IsNewUser:       false,
	})
}

func (s *Server) registerWithMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.RegisterWithMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate the registration token
	claims, err := jwt.ParseToken[*jwt.UserMagicLinkRegistrationClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if claims.Subject != jwt.UserSignatureSubjectMagicLinkRegistration {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Generate refresh token and XSRF token
	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Create a new user
	now := time.Now()
	u := &core.User{
		ID:               uuid.Must(uuid.NewV4()),
		Email:            claims.Email,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		RefreshTokenHash: hashedRefreshToken,
	}

	orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return err
	}
	hasOrganization := len(orgAccesses) > 0

	var token, xsrfToken string
	var expiration time.Duration

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		// Create the user in a transaction
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		expiration = core.TmpTokenExpiration
		if !config.Config.IsCloudEdition {
			// For self-hosted, create initial organization
			if err := s.createInitialOrganizationForSelfHosted(ctx, tx, u); err != nil {
				return err
			}
			expiration = core.TokenExpiration()
			hasOrganization = true
		}

		// Create token
		token, xsrfToken, _, _, _, err = s.createTokens(u.ID, expiration)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	cookieConfig := newCookieConfig()
	if config.Config.IsCloudEdition {
		cookieConfig.SetTmpAuthCookie(w, token, xsrfToken, config.Config.AuthDomain())
	} else {
		cookieConfig.SetAuthCookie(w, token, plainRefreshToken, xsrfToken,
			int(core.TokenExpiration().Seconds()),
			int(core.RefreshTokenExpiration.Seconds()),
			int(core.XSRFTokenExpiration.Seconds()),
			config.Config.BaseDomain)
	}

	return s.renderJSON(w, http.StatusOK, &responses.RegisterWithMagicLinkResponse{
		ExpiresAt:       strconv.FormatInt(now.Add(expiration).Unix(), 10),
		HasOrganization: hasOrganization,
	})
}

func (s *Server) requestInvitationMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.RequestInvitationMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate invitation token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](req.InvitationToken)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get invitation
	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get organization
	invitedOrg, err := s.db.Organization().Get(ctx, database.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return err
	}

	// Verify organization access in cloud edition
	if config.Config.IsCloudEdition {
		ctxSubdomain := internal.ContextSubdomain(ctx)
		hostOrg, err := s.db.Organization().Get(ctx, database.OrganizationBySubdomain(ctxSubdomain))
		if err != nil {
			return err
		}

		if invitedOrg.ID != hostOrg.ID {
			return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	// Create magic link token
	tok, err := createInvitationMagicLinkToken(c.Email)
	if err != nil {
		return err
	}

	// Build magic link URL
	url, err := buildInvitationMagicLinkURL(internal.SafeValue(invitedOrg.Subdomain), tok)
	if err != nil {
		return err
	}

	// Send magic link email
	if err := s.sendInvitationMagicLinkEmail(ctx, c.Email, "there", url); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, &responses.RequestInvitationMagicLinkResponse{
		Email: c.Email,
	})
}

func (s *Server) authenticateWithInvitationMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.AuthenticateWithInvitationMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if c.Subject != jwt.UserSignatureSubjectInvitationMagicLink {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get invitation
	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get organization
	invitedOrg, err := s.db.Organization().Get(ctx, database.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		ctxSubdomain := internal.ContextSubdomain(ctx)
		hostOrg, err := s.db.Organization().Get(ctx, database.OrganizationBySubdomain(ctxSubdomain))
		if err != nil {
			return err
		}

		if invitedOrg.ID != hostOrg.ID {
			return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = internal.SafeValue(hostOrg.Subdomain)
	}

	// Check if user exists
	exists, err := s.db.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := createMagicLinkRegistrationToken(c.Email)
		if err != nil {
			return errdefs.ErrInvalidArgument(fmt.Errorf("failed to generate registration token: %w", err))
		}

		return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithInvitationMagicLinkResponse{
			Token:     registrationToken,
			IsNewUser: true,
		})
	}

	// Get existing user
	u, err := s.db.User().Get(ctx, database.UserByEmail(c.Email))
	if err != nil {
		return err
	}

	// Create organization access
	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	// Generate token and refresh token
	now := time.Now()
	expiresAt := now.Add(core.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("failed to generate token: %w", err))
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithInvitationMagicLinkResponse{
		AuthURL:   config.Config.OrgBaseURL(orgSubdomain) + core.SaveAuthPath,
		Token:     token,
		IsNewUser: false,
	})
}

func (s *Server) registerWithInvitationMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requests.RegisterWithInvitationMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if c.Subject != jwt.UserSignatureSubjectMagicLinkRegistration {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get invitation
	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get organization
	invitedOrg, err := s.db.Organization().Get(ctx, database.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		ctxSubdomain := internal.ContextSubdomain(ctx)
		hostOrg, err := s.db.Organization().Get(ctx, database.OrganizationBySubdomain(ctxSubdomain))
		if err != nil {
			return err
		}

		if invitedOrg.ID != hostOrg.ID {
			return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = internal.SafeValue(hostOrg.Subdomain)
	}

	// Generate refresh token
	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	// Create new user
	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())
	u := &core.User{
		ID:               uuid.Must(uuid.NewV4()),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            c.Email,
		RefreshTokenHash: hashedRefreshToken,
	}

	// Create organization access
	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	// Generate token
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
			return err
		}

		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetAuthCookie(w, token, plainRefreshToken, xsrfToken,
		int(core.TokenExpiration().Seconds()),
		int(core.RefreshTokenExpiration.Seconds()),
		int(core.XSRFTokenExpiration.Seconds()),
		config.Config.OrgDomain(orgSubdomain))

	return s.renderJSON(w, http.StatusOK, &responses.RegisterWithInvitationMagicLinkResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
	})
}
