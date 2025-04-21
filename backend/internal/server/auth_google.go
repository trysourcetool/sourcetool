package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

	if err := s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
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

		return nil
	}); err != nil {
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

	if err := s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
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

		return nil
	}); err != nil {
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
