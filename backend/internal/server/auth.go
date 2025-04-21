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
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func buildSaveAuthURL(subdomain string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), core.SaveAuthPath, nil)
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

func buildLoginURL(subdomain string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("login"), nil)
}

func (s *Server) createPersonalAPIKey(ctx context.Context, tx database.Tx, u *core.User, org *core.Organization) error {
	devEnv, err := s.db.Environment().Get(ctx, database.EnvironmentByOrganizationID(org.ID), database.EnvironmentBySlug(core.EnvironmentSlugDevelopment))
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

	return tx.APIKey().Create(ctx, apiKey)
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
	hashedRefreshToken := core.HashRefreshToken(refreshTokenCookie.Value)
	u, err := s.db.User().Get(ctx, database.UserByRefreshTokenHash(hashedRefreshToken))
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	// Get current subdomain and resolve organization
	subdomain := internal.Subdomain(ctx)
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// Verify user has access to this organization
			if _, err := s.db.User().GetOrganizationAccess(ctx,
				database.UserOrganizationAccessByUserID(u.ID),
				database.UserOrganizationAccessByOrganizationSubdomain(subdomain)); err != nil {
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
	u, err := s.db.User().Get(ctx, database.UserByID(userID))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to get user: %w", err))
	}

	// Get current subdomain and verify organization access
	subdomain := internal.Subdomain(ctx)
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// Verify user has access to this organization
			if _, err := s.db.User().GetOrganizationAccess(ctx,
				database.UserOrganizationAccessByUserID(u.ID),
				database.UserOrganizationAccessByOrganizationSubdomain(subdomain)); err != nil {
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

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, u); err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to update user: %w", err))
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
	org, _, err := s.resolveOrganization(ctx, u)
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
	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, u); err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to update user: %w", err))
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, &responses.ObtainAuthTokenResponse{
		AuthURL: authURL,
		Token:   token,
	})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	u := internal.CurrentUser(ctx)

	orgAccessOpts := []database.UserOrganizationAccessQuery{
		database.UserOrganizationAccessByUserID(u.ID),
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = internal.Subdomain(ctx)
		orgAccessOpts = append(orgAccessOpts, database.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := s.db.User().GetOrganizationAccess(ctx, orgAccessOpts...)
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
