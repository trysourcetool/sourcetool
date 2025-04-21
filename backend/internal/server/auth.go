package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/google"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/internal/mail"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func buildSaveAuthURL(subdomain string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), core.SaveAuthPath, nil)
}

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

func createAuthToken(userID, xsrfToken string, expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserAuthClaims{
		UserID:    userID,
		XSRFToken: xsrfToken,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expirationTime),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	})
}

func createGoogleAuthLinkToken(flow jwt.GoogleAuthFlow, invitationOrgID uuid.UUID, hostSubdomain string) (string, error) {
	claims := &jwt.UserGoogleAuthLinkClaims{
		Flow:            flow,
		InvitationOrgID: invitationOrgID,
		HostSubdomain:   hostSubdomain,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectGoogleAuthLink,
		},
	}
	return jwt.SignToken(claims)
}

func createGoogleRegistrationToken(googleID, email, firstName, lastName string, flow jwt.GoogleAuthFlow, invitationOrgID uuid.UUID, role string) (string, error) {
	claims := &jwt.UserGoogleRegistrationClaims{
		GoogleID:        googleID,
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		Flow:            flow,
		InvitationOrgID: invitationOrgID,
		Role:            role,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectGoogleRegistration,
		},
	}
	return jwt.SignToken(claims)
}

func buildLoginURL(subdomain string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("login"), nil)
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

// hashRefreshToken creates a SHA-256 hash of a plaintext refresh token.
func hashRefreshToken(plainRefreshToken string) string {
	hash := sha256.Sum256([]byte(plainRefreshToken))
	return hex.EncodeToString(hash[:])
}

func (s *Server) createPersonalAPIKey(ctx context.Context, tx *sqlx.Tx, u *core.User, org *core.Organization) error {
	devEnv, err := s.db.GetEnvironment(ctx, postgres.EnvironmentByOrganizationID(org.ID), postgres.EnvironmentBySlug(core.EnvironmentSlugDevelopment))
	if err != nil {
		return err
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}

	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         u.ID,
		Name:           "",
		Key:            key,
	}

	return s.db.CreateAPIKey(ctx, tx, apiKey)
}

// createTokens creates a new authentication token and refresh token.
func (s *Server) createTokens(userID uuid.UUID, expiration time.Duration) (token, xsrfToken, plainRefreshToken, hashedRefreshToken string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(expiration)
	xsrfToken = uuid.Must(uuid.NewV4()).String()

	token, err = createAuthToken(userID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return "", "", "", "", time.Time{}, err
	}

	plainRefreshToken, hashedRefreshToken, err = core.GenerateRefreshToken()
	if err != nil {
		return "", "", "", "", time.Time{}, errdefs.ErrInternal(err)
	}

	return token, xsrfToken, plainRefreshToken, hashedRefreshToken, expiresAt, nil
}

// resolveOrganizationBySubdomain gets an organization by subdomain and verifies the user has access.
// Deprecated: Use getOrganizationBySubdomain instead.
func (s *Server) resolveOrganizationBySubdomain(ctx context.Context, u *core.User, subdomain string) (*core.Organization, *core.UserOrganizationAccess, error) {
	if subdomain == "" {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("subdomain cannot be empty"))
	}

	return s.getOrganizationBySubdomain(ctx, u, subdomain)
}

// validateSelfHostedOrganization checks if creating a new organization is allowed in self-hosted mode.
func (s *Server) validateSelfHostedOrganization(ctx context.Context) error {
	if !config.Config.IsCloudEdition {
		// In self-hosted mode, check if an organization already exists
		if _, err := s.db.GetOrganization(ctx); err == nil {
			return errdefs.ErrPermissionDenied(errors.New("only one organization is allowed in self-hosted edition"))
		}
	}
	return nil
}

func (s *Server) createInitialOrganizationForSelfHosted(ctx context.Context, tx *sqlx.Tx, u *core.User) error {
	if config.Config.IsCloudEdition {
		return nil
	}

	org := &core.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: nil, // Empty subdomain for non-cloud edition
	}
	if err := s.db.CreateOrganization(ctx, tx, org); err != nil {
		return err
	}

	orgAccess := &core.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: org.ID,
		Role:           core.UserOrganizationRoleAdmin,
	}
	if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
		return err
	}

	devEnv := &core.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		Name:           core.EnvironmentNameDevelopment,
		Slug:           core.EnvironmentSlugDevelopment,
		Color:          core.EnvironmentColorDevelopment,
	}
	envs := []*core.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: org.ID,
			Name:           core.EnvironmentNameProduction,
			Slug:           core.EnvironmentSlugProduction,
			Color:          core.EnvironmentColorProduction,
		},
		devEnv,
	}
	if err := s.db.BulkInsertEnvironments(ctx, tx, envs); err != nil {
		return err
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}
	apiKey := &core.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         u.ID,
		Name:           "",
		Key:            key,
	}
	if err := s.db.CreateAPIKey(ctx, tx, apiKey); err != nil {
		return err
	}

	return nil
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *Server) getUserOrganizationInfo(ctx context.Context) (*core.Organization, *core.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, internal.CurrentUser(ctx))
}

// getOrganizationBySubdomain retrieves an organization by subdomain and verifies user access.
func (s *Server) getOrganizationBySubdomain(ctx context.Context, u *core.User, subdomain string) (*core.Organization, *core.UserOrganizationAccess, error) {
	// Get organization by subdomain
	org, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	// Verify user has access to this organization
	orgAccess, err := s.db.GetUserOrganizationAccess(ctx,
		postgres.UserOrganizationAccessByOrganizationID(org.ID),
		postgres.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *Server) getOrganizationInfo(ctx context.Context, u *core.User) (*core.Organization, *core.UserOrganizationAccess, error) {
	if u == nil {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("user cannot be nil"))
	}

	subdomain := internal.Subdomain(ctx)
	isCloudWithSubdomain := config.Config.IsCloudEdition && subdomain != "" && subdomain != "auth"

	// Different strategies for cloud vs. self-hosted or auth subdomain
	if isCloudWithSubdomain {
		return s.getOrganizationBySubdomain(ctx, u, subdomain)
	}

	return s.getDefaultOrganizationForUser(ctx, u)
}

// (typically the most recently created one).
func (s *Server) getDefaultOrganizationForUser(ctx context.Context, u *core.User) (*core.Organization, *core.UserOrganizationAccess, error) {
	// Get user's organization access
	orgAccess, err := s.db.GetUserOrganizationAccess(ctx,
		postgres.UserOrganizationAccessByUserID(u.ID),
		postgres.UserOrganizationAccessOrderBy("created_at DESC"))
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
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

	if err := s.mail.Send(ctx, mail.MailInput{
		From:    "Sourcetool Team",
		To:      []string{email},
		Subject: subject,
		Body:    content,
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

	if err := s.mail.Send(ctx, mail.MailInput{
		From:    "Sourcetool Team",
		To:      []string{email},
		Subject: subject,
		Body:    content,
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

	if err := s.mail.Send(ctx, mail.MailInput{
		From:    "Sourcetool Team",
		To:      []string{email},
		Subject: subject,
		Body:    content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Server) sendMultipleOrganizationsLoginEmail(ctx context.Context, email, firstName string, loginURLs []string) error {
	subject := "Choose your Sourcetool organization to log in"

	urlList := ""
	for _, url := range loginURLs {
		urlList += url + "\n"
	}

	content := fmt.Sprintf(`Hi %s,

Your email, %s, is associated with multiple Sourcetool organizations. You may log in to each one by clicking its login link below:

%s

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, email, urlList)

	if err := s.mail.Send(ctx, mail.MailInput{
		From:    "Sourcetool Team",
		To:      []string{email},
		Subject: subject,
		Body:    content,
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
	exists, err := s.db.IsUserEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}

	var firstName string
	isNewUser := !exists

	// Handle Cloud Edition with subdomain
	if config.Config.IsCloudEdition {
		subdomain := internal.Subdomain(ctx)
		if subdomain != "" && subdomain != "auth" {
			// Get organization by subdomain
			org, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
			if err != nil {
				return err
			}

			if exists {
				// For existing users, check if they have access to this organization
				u, err := s.db.GetUser(ctx, postgres.UserByEmail(req.Email))
				if err != nil {
					return err
				}

				_, err = s.db.GetUserOrganizationAccess(ctx,
					postgres.UserOrganizationAccessByUserID(u.ID),
					postgres.UserOrganizationAccessByOrganizationID(org.ID))
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
		u, err := s.db.GetUser(ctx, postgres.UserByEmail(req.Email))
		if err != nil {
			return err
		}
		firstName = u.FirstName

		// Get user's organization access information
		orgAccesses, err := s.db.ListUserOrganizationAccesses(ctx, postgres.UserOrganizationAccessByUserID(u.ID))
		if err != nil {
			return err
		}

		// Cloud edition specific handling for multiple organizations
		if config.Config.IsCloudEdition && len(orgAccesses) > 1 {
			// Handle multiple organizations
			loginURLs := make([]string, 0, len(orgAccesses))
			for _, access := range orgAccesses {
				org, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(access.OrganizationID))
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
		subdomain = internal.Subdomain(ctx)
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
	exists, err := s.db.IsUserEmailExists(ctx, c.Email)
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
	u, err := s.db.GetUser(ctx, postgres.UserByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get user's organization access information
	orgAccesses, err := s.db.ListUserOrganizationAccesses(ctx, postgres.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return err
	}

	// Handle organization subdomain logic
	subdomain := internal.Subdomain(ctx)
	var orgAccess *core.UserOrganizationAccess
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain, resolve org and access
			_, orgAccess, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain
			if len(orgAccesses) == 0 {
				// No organization - sign in as a user not associated with any organization
			} else if len(orgAccesses) == 1 {
				// Single organization - redirect to it
				orgAccess = orgAccesses[0]
				org, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(orgAccess.OrganizationID))
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
		_, err = s.db.GetOrganization(ctx, postgres.OrganizationByID(orgAccess.OrganizationID))
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

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save changes
	if err = s.db.UpdateUser(ctx, tx, u); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
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

	orgAccesses, err := s.db.ListUserOrganizationAccesses(ctx, postgres.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return err
	}
	hasOrganization := len(orgAccesses) > 0

	var token, xsrfToken string
	var expiration time.Duration

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create the user in a transaction
	if err := s.db.CreateUser(ctx, tx, u); err != nil {
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

	if err = tx.Commit(); err != nil {
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
	userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get organization
	invitedOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return err
	}

	// Verify organization access in cloud edition
	if config.Config.IsCloudEdition {
		subdomain := internal.Subdomain(ctx)
		hostOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
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
	userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get organization
	invitedOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := internal.Subdomain(ctx)
		hostOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
		if err != nil {
			return err
		}

		if invitedOrg.ID != hostOrg.ID {
			return errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = internal.SafeValue(hostOrg.Subdomain)
	}

	// Check if user exists
	exists, err := s.db.IsUserEmailExists(ctx, c.Email)
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
	u, err := s.db.GetUser(ctx, postgres.UserByEmail(c.Email))
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

	// Save changes
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = s.db.DeleteUserInvitation(ctx, tx, userInvitation); err != nil {
		return err
	}

	if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
		return err
	}

	if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
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
	userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(c.Email))
	if err != nil {
		return err
	}

	// Get organization
	invitedOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := internal.Subdomain(ctx)
		hostOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
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

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save changes
	if err := s.db.DeleteUserInvitation(ctx, tx, userInvitation); err != nil {
		return err
	}

	if err := s.db.CreateUser(ctx, tx, u); err != nil {
		return err
	}

	if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
		return err
	}

	if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
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

func (s *Server) requestGoogleAuthLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var hostSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := internal.Subdomain(ctx)
		if subdomain != "auth" {
			hostSubdomain = subdomain
		}
	}

	stateToken, err := createGoogleAuthLinkToken(
		jwt.GoogleAuthFlowStandard,
		uuid.Nil,
		hostSubdomain,
	)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	googleOAuthClient := google.NewOAuthClient()
	url, err := googleOAuthClient.GetGoogleAuthCodeURL(ctx, stateToken)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	return s.renderJSON(w, http.StatusOK, &responses.RequestGoogleAuthLinkResponse{
		AuthURL: url,
	})
}

func (s *Server) authenticateWithGoogle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.AuthenticateWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate state token
	stateClaims, err := jwt.ParseToken[*jwt.UserGoogleAuthLinkClaims](req.State)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if stateClaims.Subject != jwt.UserSignatureSubjectGoogleAuthLink {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get Google token and user info
	googleOAuthClient := google.NewOAuthClient()
	tok, err := googleOAuthClient.GetGoogleToken(ctx, req.Code)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	userInfo, err := googleOAuthClient.GetGoogleUserInfo(ctx, tok)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	// In staging environment, only allow @trysourcetool.com email addresses
	if config.Config.Env == config.EnvStaging && !strings.HasSuffix(userInfo.Email, "@trysourcetool.com") {
		return errdefs.ErrPermissionDenied(errors.New("access restricted in staging environment"))
	}

	// Check if user exists
	exists, err := s.db.IsUserEmailExists(ctx, userInfo.Email)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	if !exists {
		if !config.Config.IsCloudEdition && stateClaims.Flow == jwt.GoogleAuthFlowStandard {
			if err := s.validateSelfHostedOrganization(ctx); err != nil {
				return errdefs.ErrInternal(err)
			}
		}

		var role string
		if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
			// Verify invitation exists
			userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(userInfo.Email), postgres.UserInvitationByOrganizationID(stateClaims.InvitationOrgID))
			if err != nil {
				return errdefs.ErrInvalidArgument(errors.New("invalid invitation"))
			}
			role = userInvitation.Role.String()
		}

		// Generate registration token with flow info
		registrationToken, err := createGoogleRegistrationToken(
			userInfo.ID,
			userInfo.Email,
			userInfo.GivenName,
			userInfo.FamilyName,
			stateClaims.Flow,
			stateClaims.InvitationOrgID,
			role,
		)
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to create registration token: %w", err))
		}

		if stateClaims.Flow != jwt.GoogleAuthFlowInvitation {
			cookieConfig := newCookieConfig()
			xsrfToken := uuid.Must(uuid.NewV4()).String()
			cookieConfig.SetTmpAuthCookie(w, registrationToken, xsrfToken, config.Config.AuthDomain())
		}

		return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithGoogleResponse{
			Token:           registrationToken,
			IsNewUser:       true,
			HasOrganization: stateClaims.Flow == jwt.GoogleAuthFlowInvitation,
			FirstName:       userInfo.GivenName,
			LastName:        userInfo.FamilyName,
		})
	}

	// For existing users
	u, err := s.db.GetUser(ctx, postgres.UserByEmail(userInfo.Email))
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	needsGoogleIDUpdate := u.GoogleID == ""

	var org *core.Organization
	var orgAccess *core.UserOrganizationAccess
	var orgSubdomain string

	if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
		// Handle invitation flow for existing users
		invitedOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(stateClaims.InvitationOrgID))
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to get invited organization: %w", err))
		}

		userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(userInfo.Email), postgres.UserInvitationByOrganizationID(stateClaims.InvitationOrgID))
		if err != nil {
			return errdefs.ErrInvalidArgument(errors.New("invalid invitation"))
		}

		orgAccess = &core.UserOrganizationAccess{
			ID:             uuid.Must(uuid.NewV4()),
			UserID:         u.ID,
			OrganizationID: invitedOrg.ID,
			Role:           userInvitation.Role,
		}
		org = invitedOrg
		orgSubdomain = internal.SafeValue(invitedOrg.Subdomain)
	} else {
		// Standard flow - get user's organization info
		// Get all organization accesses for the user
		orgAccesses, err := s.db.ListUserOrganizationAccesses(ctx, postgres.UserOrganizationAccessByUserID(u.ID))
		if err != nil {
			return errdefs.ErrInternal(err)
		}

		if config.Config.IsCloudEdition {
			if len(orgAccesses) > 1 {
				hostSubdomain := stateClaims.HostSubdomain
				if hostSubdomain == "" {
					// Handle multiple organizations by sending email with login URLs
					loginURLs := make([]string, 0, len(orgAccesses))
					for _, access := range orgAccesses {
						org, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(access.OrganizationID))
						if err != nil {
							return errdefs.ErrInternal(err)
						}

						url, err := buildLoginURL(internal.SafeValue(org.Subdomain))
						if err != nil {
							return errdefs.ErrInternal(err)
						}
						loginURLs = append(loginURLs, url)
					}

					// Send email with multiple organization links
					if err := s.sendMultipleOrganizationsLoginEmail(ctx, u.Email, u.FirstName, loginURLs); err != nil {
						return errdefs.ErrInternal(err)
					}

					return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithGoogleResponse{
						IsNewUser:                false,
						HasOrganization:          true,
						HasMultipleOrganizations: true,
					})
				} else {
					org, err = s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(hostSubdomain))
					if err != nil {
						return errdefs.ErrInternal(err)
					}
					orgAccess, err = s.db.GetUserOrganizationAccess(ctx, postgres.UserOrganizationAccessByUserID(u.ID), postgres.UserOrganizationAccessByOrganizationID(org.ID))
					if err != nil {
						return err
					}
					orgSubdomain = internal.SafeValue(org.Subdomain)
				}
			} else {
				// Single organization case
				orgAccess = orgAccesses[0]

				org, err = s.db.GetOrganization(ctx, postgres.OrganizationByID(orgAccess.OrganizationID))
				if err != nil {
					return errdefs.ErrInternal(err)
				}
				orgSubdomain = internal.SafeValue(org.Subdomain)
			}
		} else {
			// Self-hosted mode
			orgAccess = orgAccesses[0]
			org, err = s.db.GetOrganization(ctx, postgres.OrganizationByID(orgAccess.OrganizationID))
			if err != nil {
				return errdefs.ErrInternal(err)
			}
		}
	}

	// Generate temporary auth tokens
	token, xsrfToken, _, hashedRefreshToken, _, err := s.createTokens(u.ID, core.TmpTokenExpiration)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	u.RefreshTokenHash = hashedRefreshToken
	if needsGoogleIDUpdate {
		u.GoogleID = userInfo.ID
	}

	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
		// For invitation flow, create org access and delete invitation
		userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(userInfo.Email), postgres.UserInvitationByOrganizationID(stateClaims.InvitationOrgID))
		if err != nil {
			return err
		}
		if err := s.db.DeleteUserInvitation(ctx, tx, userInvitation); err != nil {
			return err
		}
		if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
			return err
		}
		if err := s.createPersonalAPIKey(ctx, tx, u, org); err != nil {
			return err
		}
	}

	if err := s.db.UpdateUser(ctx, tx, u); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if orgAccess == nil && stateClaims.Flow != jwt.GoogleAuthFlowInvitation {
		cookieConfig := newCookieConfig()
		cookieConfig.SetTmpAuthCookie(w, token, xsrfToken, config.Config.AuthDomain())
	}

	return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithGoogleResponse{
		AuthURL:         authURL,
		Token:           token,
		HasOrganization: orgAccess != nil,
		IsNewUser:       false,
	})
}

func (s *Server) registerWithGoogle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.RegisterWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate registration token
	claims, err := jwt.ParseToken[*jwt.UserGoogleRegistrationClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid registration token: %w", err))
	}
	if claims.Subject != jwt.UserSignatureSubjectGoogleRegistration {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject for google registration"))
	}

	// Check if user already exists
	exists, err := s.db.IsUserEmailExists(ctx, claims.Email)
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to check user existence: %w", err))
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(fmt.Errorf("user with email %s already exists", claims.Email))
	}

	_, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to generate refresh token: %w", err))
	}

	tokenExpiration := core.TokenExpiration()
	u := &core.User{
		ID:               uuid.Must(uuid.NewV4()),
		Email:            claims.Email,
		FirstName:        claims.FirstName,
		LastName:         claims.LastName,
		RefreshTokenHash: hashedRefreshToken,
		GoogleID:         claims.GoogleID,
	}

	var token, xsrfToken string
	var orgSubdomain string
	var authURL string
	var hasOrganization bool

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.db.CreateUser(ctx, tx, u); err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to create user: %w", err))
	}

	if claims.Flow == jwt.GoogleAuthFlowInvitation {
		invitedOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(claims.InvitationOrgID))
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to get invited organization: %w", err))
		}

		userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(claims.Email), postgres.UserInvitationByOrganizationID(claims.InvitationOrgID))
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to get invitation: %w", err))
		}

		orgAccess := &core.UserOrganizationAccess{
			ID:             uuid.Must(uuid.NewV4()),
			UserID:         u.ID,
			OrganizationID: claims.InvitationOrgID,
			Role:           core.UserOrganizationRoleFromString(claims.Role),
		}

		if err := s.db.DeleteUserInvitation(ctx, tx, userInvitation); err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to delete invitation: %w", err))
		}

		if err := s.db.CreateUserOrganizationAccess(ctx, tx, orgAccess); err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to create organization access: %w", err))
		}

		if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to create personal API key: %w", err))
		}

		orgSubdomain = internal.SafeValue(invitedOrg.Subdomain)
		hasOrganization = true
	} else {
		if !config.Config.IsCloudEdition {
			if err := s.createInitialOrganizationForSelfHosted(ctx, tx, u); err != nil {
				return errdefs.ErrInternal(fmt.Errorf("failed to create initial organization: %w", err))
			}
			hasOrganization = true
		}
	}

	token, xsrfToken, _, _, _, err = s.createTokens(u.ID, tokenExpiration)
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to create auth token: %w", err))
	}

	if hasOrganization {
		authURL, err = buildSaveAuthURL(orgSubdomain)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetTmpAuthCookie(w, token, xsrfToken, config.Config.AuthDomain())

	return s.renderJSON(w, http.StatusOK, &responses.RegisterWithGoogleResponse{
		Token:           token,
		AuthURL:         authURL,
		HasOrganization: hasOrganization,
	})
}

func (s *Server) requestInvitationGoogleAuthLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.RequestInvitationGoogleAuthLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	c, err := jwt.ParseToken[*jwt.UserEmailClaims](req.InvitationToken)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid invitation token: %w", err))
	}
	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return errdefs.ErrInvalidArgument(errors.New("invalid jwt subject for invitation"))
	}

	userInvitation, err := s.db.GetUserInvitation(ctx, postgres.UserInvitationByEmail(c.Email))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to retrieve invitation: %w", err))
	}

	invitedOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to retrieve invited organization: %w", err))
	}

	if config.Config.IsCloudEdition {
		subdomain := internal.Subdomain(ctx)
		if subdomain == "" || subdomain == "auth" {
			return errdefs.ErrInvalidArgument(errors.New("invitation must be accessed via organization subdomain"))
		}
		hostOrg, err := s.db.GetOrganization(ctx, postgres.OrganizationBySubdomain(subdomain))
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to retrieve host organization: %w", err))
		}
		if invitedOrg.ID != hostOrg.ID {
			return errdefs.ErrUnauthenticated(errors.New("invitation organization mismatch"))
		}
	}

	var hostSubdomain string
	if config.Config.IsCloudEdition {
		hostSubdomain = internal.Subdomain(ctx)
	}

	stateToken, err := createGoogleAuthLinkToken(
		jwt.GoogleAuthFlowInvitation,
		invitedOrg.ID,
		hostSubdomain,
	)
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to create state token: %w", err))
	}

	googleOAuthClient := google.NewOAuthClient()
	url, err := googleOAuthClient.GetGoogleAuthCodeURL(ctx, stateToken)
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to get google auth code url: %w", err))
	}

	return s.renderJSON(w, http.StatusOK, &responses.RequestInvitationGoogleAuthLinkResponse{
		AuthURL: url,
	})
}

func (s *Server) refreshToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
	if xsrfTokenHeader == "" {
		return errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token"))
	}

	xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	// Validate XSRF token consistency
	if xsrfTokenCookie.Value != xsrfTokenHeader {
		return errdefs.ErrUnauthenticated(errors.New("invalid xsrf token"))
	}

	// Get user by refresh token
	hashedRefreshToken := hashRefreshToken(refreshTokenCookie.Value)
	u, err := s.db.GetUser(ctx, postgres.UserByRefreshTokenHash(hashedRefreshToken))
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	// Get current subdomain and resolve organization
	subdomain := internal.Subdomain(ctx)
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// Verify user has access to this organization
			_, _, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain, use default
			orgSubdomain = "auth"
		}
	} else {
		// For self-hosted, no specific subdomain needed
		orgSubdomain = ""
	}

	// Generate token and set expiration
	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetAuthCookie(w, token, refreshTokenCookie.Value, xsrfToken,
		int(core.TokenExpiration().Seconds()),
		int(core.RefreshTokenExpiration.Seconds()),
		int(core.XSRFTokenExpiration.Seconds()),
		config.Config.OrgDomain(orgSubdomain))

	return s.renderJSON(w, http.StatusOK, &responses.RefreshTokenResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
	})
}

func (s *Server) saveAuth(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req requests.SaveAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserAuthClaims](req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid token: %w", err))
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid user id: %w", err))
	}

	// Get user by ID
	u, err := s.db.GetUser(ctx, postgres.UserByID(userID))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to get user: %w", err))
	}

	// Get current subdomain and verify organization access
	subdomain := internal.Subdomain(ctx)
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain, verify user has access
			_, _, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain, use default
			orgSubdomain = "auth"
		}
	}

	// Generate token and refresh token
	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	// Update user's refresh token
	u.RefreshTokenHash = hashedRefreshToken

	tx, err := s.db.Beginx()
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback()

	if err := s.db.UpdateUser(ctx, tx, u); err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to update user: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetAuthCookie(w, token, plainRefreshToken, xsrfToken,
		int(core.TokenExpiration().Seconds()),
		int(core.RefreshTokenExpiration.Seconds()),
		int(core.XSRFTokenExpiration.Seconds()),
		config.Config.OrgDomain(orgSubdomain))

	return s.renderJSON(w, http.StatusOK, &responses.SaveAuthResponse{
		ExpiresAt:   strconv.FormatInt(expiresAt.Unix(), 10),
		RedirectURL: config.Config.OrgBaseURL(orgSubdomain),
	})
}

func (s *Server) obtainAuthToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	// Get current user from context
	u := internal.CurrentUser(ctx)
	if u == nil {
		return errdefs.ErrUnauthenticated(errors.New("no user in context"))
	}

	// Get user's organization info
	org, _, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return err
	}

	// Generate temporary token
	now := time.Now()
	expiresAt := now.Add(core.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return err
	}

	// Build auth URL with organization subdomain
	authURL, err := buildSaveAuthURL(internal.SafeValue(org.Subdomain))
	if err != nil {
		return err
	}

	// Update user
	tx, err := s.db.Beginx()
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback()

	if err = s.db.UpdateUser(ctx, tx, u); err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to update user: %w", err))
	}

	if err = tx.Commit(); err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to commit transaction: %w", err))
	}

	return s.renderJSON(w, http.StatusOK, &responses.ObtainAuthTokenResponse{
		AuthURL: authURL,
		Token:   token,
	})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	u := internal.CurrentUser(ctx)

	orgAccessOpts := []postgres.UserOrganizationAccessQuery{
		postgres.UserOrganizationAccessByUserID(u.ID),
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = internal.Subdomain(ctx)
		orgAccessOpts = append(orgAccessOpts, postgres.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := s.db.GetUserOrganizationAccess(ctx, orgAccessOpts...)
	if err != nil {
		return err
	}

	cookieConfig := newCookieConfig()
	cookieConfig.DeleteAuthCookie(w, r, config.Config.OrgDomain(subdomain))

	return s.renderJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully logged out",
	})
}
