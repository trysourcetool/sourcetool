package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/google"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/internal/mail"
	"github.com/trysourcetool/sourcetool/backend/internal/server/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) requestGoogleAuthLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var hostSubdomain string
	if config.Config.IsCloudEdition {
		ctxSubdomain := internal.ContextSubdomain(ctx)
		if ctxSubdomain != "auth" {
			hostSubdomain = ctxSubdomain
		}
	}

	stateToken, err := jwt.SignGoogleAuthLinkToken(
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
	stateClaims, err := jwt.ParseGoogleAuthLinkClaims(req.State)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
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
	exists, err := s.db.User().IsEmailExists(ctx, userInfo.Email)
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
			userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(userInfo.Email), database.UserInvitationByOrganizationID(stateClaims.InvitationOrgID))
			if err != nil {
				return errdefs.ErrInvalidArgument(errors.New("invalid invitation"))
			}
			role = userInvitation.Role.String()
		}

		// Generate registration token with flow info
		registrationToken, err := jwt.SignGoogleRegistrationToken(
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
	u, err := s.db.User().Get(ctx, database.UserByEmail(userInfo.Email))
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	needsGoogleIDUpdate := u.GoogleID == ""

	var org *core.Organization
	var orgAccess *core.UserOrganizationAccess
	var orgSubdomain string

	if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
		// Handle invitation flow for existing users
		invitedOrg, err := s.db.Organization().Get(ctx, database.OrganizationByID(stateClaims.InvitationOrgID))
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to get invited organization: %w", err))
		}

		userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(userInfo.Email), database.UserInvitationByOrganizationID(stateClaims.InvitationOrgID))
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
		orgSubdomain = internal.StringValue(invitedOrg.Subdomain)
	} else {
		// Standard flow - get user's organization info
		// Get all organization accesses for the user
		orgAccesses, err := s.db.User().ListOrganizationAccesses(ctx, database.UserOrganizationAccessByUserID(u.ID))
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
						org, err := s.db.Organization().Get(ctx, database.OrganizationByID(access.OrganizationID))
						if err != nil {
							return errdefs.ErrInternal(err)
						}

						url, err := buildLoginURL(internal.StringValue(org.Subdomain))
						if err != nil {
							return errdefs.ErrInternal(err)
						}
						loginURLs = append(loginURLs, url)
					}

					// Send email with multiple organization links
					if err := mail.SendMultipleOrganizationsLoginEmail(ctx, u.Email, u.FirstName, loginURLs); err != nil {
						return errdefs.ErrInternal(err)
					}

					return s.renderJSON(w, http.StatusOK, &responses.AuthenticateWithGoogleResponse{
						IsNewUser:                false,
						HasOrganization:          true,
						HasMultipleOrganizations: true,
					})
				} else {
					org, err = s.db.Organization().Get(ctx, database.OrganizationBySubdomain(hostSubdomain))
					if err != nil {
						return errdefs.ErrInternal(err)
					}
					orgAccess, err = s.db.User().GetOrganizationAccess(ctx, database.UserOrganizationAccessByUserID(u.ID), database.UserOrganizationAccessByOrganizationID(org.ID))
					if err != nil {
						return err
					}
					orgSubdomain = internal.StringValue(org.Subdomain)
				}
			} else {
				// Single organization case
				orgAccess = orgAccesses[0]

				org, err = s.db.Organization().Get(ctx, database.OrganizationByID(orgAccess.OrganizationID))
				if err != nil {
					return errdefs.ErrInternal(err)
				}
				orgSubdomain = internal.StringValue(org.Subdomain)
			}
		} else {
			// Self-hosted mode
			orgAccess = orgAccesses[0]
			org, err = s.db.Organization().Get(ctx, database.OrganizationByID(orgAccess.OrganizationID))
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

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
			// For invitation flow, create org access and delete invitation
			userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(userInfo.Email), database.UserInvitationByOrganizationID(stateClaims.InvitationOrgID))
			if err != nil {
				return err
			}
			if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
				return err
			}
			if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
				return err
			}
			if err := s.createPersonalAPIKey(ctx, tx, u, org); err != nil {
				return err
			}
		}

		if err := tx.User().Update(ctx, u); err != nil {
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
	claims, err := jwt.ParseGoogleRegistrationClaims(req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid registration token: %w", err))
	}

	// Check if user already exists
	exists, err := s.db.User().IsEmailExists(ctx, claims.Subject)
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to check user existence: %w", err))
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(fmt.Errorf("user with email %s already exists", claims.Subject))
	}

	_, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to generate refresh token: %w", err))
	}

	tokenExpiration := core.TokenExpiration()
	u := &core.User{
		ID:               uuid.Must(uuid.NewV4()),
		Email:            claims.Subject,
		FirstName:        claims.FirstName,
		LastName:         claims.LastName,
		RefreshTokenHash: hashedRefreshToken,
		GoogleID:         claims.GoogleID,
	}

	var token, xsrfToken string
	var orgSubdomain string
	var authURL string
	var hasOrganization bool

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to create user: %w", err))
		}

		if claims.Flow == jwt.GoogleAuthFlowInvitation {
			invitedOrg, err := s.db.Organization().Get(ctx, database.OrganizationByID(claims.InvitationOrgID))
			if err != nil {
				return errdefs.ErrInternal(fmt.Errorf("failed to get invited organization: %w", err))
			}

			userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(claims.Subject), database.UserInvitationByOrganizationID(claims.InvitationOrgID))
			if err != nil {
				return errdefs.ErrInternal(fmt.Errorf("failed to get invitation: %w", err))
			}

			orgAccess := &core.UserOrganizationAccess{
				ID:             uuid.Must(uuid.NewV4()),
				UserID:         u.ID,
				OrganizationID: claims.InvitationOrgID,
				Role:           core.UserOrganizationRoleFromString(claims.Role),
			}

			if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
				return errdefs.ErrInternal(fmt.Errorf("failed to delete invitation: %w", err))
			}

			if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
				return errdefs.ErrInternal(fmt.Errorf("failed to create organization access: %w", err))
			}

			if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
				return errdefs.ErrInternal(fmt.Errorf("failed to create personal API key: %w", err))
			}

			orgSubdomain = internal.StringValue(invitedOrg.Subdomain)
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

	c, err := jwt.ParseInvitationClaims(req.InvitationToken)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid invitation token: %w", err))
	}

	userInvitation, err := s.db.User().GetInvitation(ctx, database.UserInvitationByEmail(c.Subject))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to retrieve invitation: %w", err))
	}

	invitedOrg, err := s.db.Organization().Get(ctx, database.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to retrieve invited organization: %w", err))
	}

	if config.Config.IsCloudEdition {
		ctxSubdomain := internal.ContextSubdomain(ctx)
		if ctxSubdomain == "" || ctxSubdomain == "auth" {
			return errdefs.ErrInvalidArgument(errors.New("invitation must be accessed via organization subdomain"))
		}
		hostOrg, err := s.db.Organization().Get(ctx, database.OrganizationBySubdomain(ctxSubdomain))
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to retrieve host organization: %w", err))
		}
		if invitedOrg.ID != hostOrg.ID {
			return errdefs.ErrUnauthenticated(errors.New("invitation organization mismatch"))
		}
	}

	var ctxSubdomain string
	if config.Config.IsCloudEdition {
		ctxSubdomain = internal.ContextSubdomain(ctx)
	}

	stateToken, err := jwt.SignGoogleAuthLinkToken(
		jwt.GoogleAuthFlowInvitation,
		invitedOrg.ID,
		ctxSubdomain,
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
