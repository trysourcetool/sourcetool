package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/auth"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

type Service interface {
	// Passwordless Authentication methods
	RequestMagicLink(context.Context, dto.RequestMagicLinkInput) (*dto.RequestMagicLinkOutput, error)
	AuthenticateWithMagicLink(context.Context, dto.AuthenticateWithMagicLinkInput) (*dto.AuthenticateWithMagicLinkOutput, error)
	RegisterWithMagicLink(context.Context, dto.RegisterWithMagicLinkInput) (*dto.RegisterWithMagicLinkOutput, error)
	RequestInvitationMagicLink(context.Context, dto.RequestInvitationMagicLinkInput) (*dto.RequestInvitationMagicLinkOutput, error)
	AuthenticateWithInvitationMagicLink(context.Context, dto.AuthenticateWithInvitationMagicLinkInput) (*dto.AuthenticateWithInvitationMagicLinkOutput, error)
	RegisterWithInvitationMagicLink(context.Context, dto.RegisterWithInvitationMagicLinkInput) (*dto.RegisterWithInvitationMagicLinkOutput, error)

	// Google Authentication methods
	RequestGoogleAuthLink(context.Context) (*dto.RequestGoogleAuthLinkOutput, error)
	AuthenticateWithGoogle(context.Context, dto.AuthenticateWithGoogleInput) (*dto.AuthenticateWithGoogleOutput, error)
	RegisterWithGoogle(context.Context, dto.RegisterWithGoogleInput) (*dto.RegisterWithGoogleOutput, error)
	RequestInvitationGoogleAuthLink(context.Context, dto.RequestInvitationGoogleAuthLinkInput) (*dto.RequestInvitationGoogleAuthLinkOutput, error)

	// Authentication methods
	Logout(context.Context) (*dto.LogoutOutput, error)
	Save(context.Context, dto.SaveAuthInput) (*dto.SaveAuthOutput, error)
	RefreshToken(context.Context, dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error)
	ObtainAuthToken(context.Context) (*dto.ObtainAuthTokenOutput, error)
}

type ServiceCE struct {
	*port.Dependencies
}

func NewServiceCE(d *port.Dependencies) *ServiceCE {
	return &ServiceCE{Dependencies: d}
}

func (s *ServiceCE) RequestMagicLink(ctx context.Context, in dto.RequestMagicLinkInput) (*dto.RequestMagicLinkOutput, error) {
	// Check if email exists
	exists, err := s.Repository.User().IsEmailExists(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	var firstName string
	isNewUser := !exists

	// Handle Cloud Edition with subdomain
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		if subdomain != "" && subdomain != "auth" {
			// Get organization by subdomain
			org, err := s.Repository.Organization().Get(ctx, organization.BySubdomain(subdomain))
			if err != nil {
				return nil, err
			}

			if exists {
				// For existing users, check if they have access to this organization
				u, err := s.Repository.User().Get(ctx, user.ByEmail(in.Email))
				if err != nil {
					return nil, err
				}

				_, err = s.Repository.User().GetOrganizationAccess(ctx,
					user.OrganizationAccessByUserID(u.ID),
					user.OrganizationAccessByOrganizationID(org.ID))
				if err != nil {
					return nil, errdefs.ErrUnauthenticated(errors.New("user does not have access to this organization"))
				}
			} else {
				// For new users, registration is only allowed through invitations
				return nil, errdefs.ErrPermissionDenied(errors.New("registration is only allowed through invitations"))
			}
		}
	}

	if exists {
		// Get user by email for existing users
		u, err := s.Repository.User().Get(ctx, user.ByEmail(in.Email))
		if err != nil {
			return nil, err
		}
		firstName = u.FirstName

		// Get user's organization access information
		orgAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByUserID(u.ID))
		if err != nil {
			return nil, err
		}

		// Cloud edition specific handling for multiple organizations
		if config.Config.IsCloudEdition && len(orgAccesses) > 1 {
			// Handle multiple organizations
			loginURLs := make([]string, 0, len(orgAccesses))
			for _, access := range orgAccesses {
				org, err := s.Repository.Organization().Get(ctx, organization.ByID(access.OrganizationID))
				if err != nil {
					return nil, err
				}

				// Create org-specific magic link
				tok, err := createMagicLinkToken(in.Email)
				if err != nil {
					return nil, err
				}

				url, err := buildMagicLinkURL(ptrconv.SafeValue(org.Subdomain), tok)
				if err != nil {
					return nil, err
				}
				loginURLs = append(loginURLs, url)
			}

			if err := s.sendMultipleOrganizationsMagicLinkEmail(ctx, in.Email, firstName, loginURLs); err != nil {
				return nil, err
			}

			return &dto.RequestMagicLinkOutput{
				Email: in.Email,
				IsNew: false,
			}, nil
		}
	} else {
		// For new users, generate a temporary ID that will be verified/used later
		firstName = "there" // Default greeting

		// For self-hosted mode, check if creating an organization is allowed
		if !config.Config.IsCloudEdition {
			// Check if an organization already exists in self-hosted mode
			if err := s.validateSelfHostedOrganization(ctx); err != nil {
				return nil, err
			}
		}
	}

	// Determine subdomain context based on edition
	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
	}

	// Create token for magic link authentication
	tok, err := createMagicLinkToken(in.Email)
	if err != nil {
		return nil, err
	}

	// Build magic link URL
	url, err := buildMagicLinkURL(subdomain, tok)
	if err != nil {
		return nil, err
	}

	// Send magic link email
	if err := s.sendMagicLinkEmail(ctx, in.Email, firstName, url); err != nil {
		return nil, err
	}

	return &dto.RequestMagicLinkOutput{
		Email: in.Email,
		IsNew: isNewUser,
	}, nil
}

func (s *ServiceCE) AuthenticateWithMagicLink(ctx context.Context, in dto.AuthenticateWithMagicLinkInput) (*dto.AuthenticateWithMagicLinkOutput, error) {
	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectMagicLink {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Check if user exists
	exists, err := s.Repository.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := createMagicLinkRegistrationToken(c.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to generate registration token: %w", err)
		}

		return &dto.AuthenticateWithMagicLinkOutput{
			Token:           registrationToken,
			IsNewUser:       true,
			HasOrganization: false,
		}, nil
	}

	// Get existing user
	u, err := s.Repository.User().Get(ctx, user.ByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get user's organization access information
	orgAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}

	// Handle organization subdomain logic
	subdomain := ctxutil.Subdomain(ctx)
	var orgAccess *user.UserOrganizationAccess
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain, resolve org and access
			_, orgAccess, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return nil, err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain
			if len(orgAccesses) == 0 {
				// No organization - sign in as a user not associated with any organization
			} else if len(orgAccesses) == 1 {
				// Single organization - redirect to it
				orgAccess = orgAccesses[0]
				org, err := s.Repository.Organization().Get(ctx, organization.ByID(orgAccess.OrganizationID))
				if err != nil {
					return nil, err
				}
				orgSubdomain = ptrconv.SafeValue(org.Subdomain)
			} else {
				return nil, errdefs.ErrUserMultipleOrganizations(errors.New("user has multiple organizations"))
			}
		}
	} else {
		// Self-hosted mode has only one organization
		orgAccess = orgAccesses[0]
		_, err = s.Repository.Organization().Get(ctx, organization.ByID(orgAccess.OrganizationID))
		if err != nil {
			return nil, err
		}
	}

	// Create token, refresh token, etc.
	token, xsrfToken, plainRefreshToken, hashedRefreshToken, _, err := s.createTokens(
		u.ID, auth.TmpTokenExpiration)
	if err != nil {
		return nil, err
	}

	// Update user with new refresh token
	u.RefreshTokenHash = hashedRefreshToken
	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err = s.Repository.RunTransaction(func(tx port.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.AuthenticateWithMagicLinkOutput{
		AuthURL:         authURL,
		Token:           token,
		HasOrganization: orgAccess != nil,
		RefreshToken:    plainRefreshToken,
		XSRFToken:       xsrfToken,
		Domain:          config.Config.OrgDomain(orgSubdomain),
		IsNewUser:       false,
	}, nil
}

func (s *ServiceCE) RegisterWithMagicLink(ctx context.Context, in dto.RegisterWithMagicLinkInput) (*dto.RegisterWithMagicLinkOutput, error) {
	// Parse and validate the registration token
	claims, err := jwt.ParseToken[*jwt.UserMagicLinkRegistrationClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if claims.Subject != jwt.UserSignatureSubjectMagicLinkRegistration {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Generate refresh token and XSRF token
	plainRefreshToken, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create a new user
	now := time.Now()
	u := &user.User{
		ID:               uuid.Must(uuid.NewV4()),
		Email:            claims.Email,
		FirstName:        in.FirstName,
		LastName:         in.LastName,
		RefreshTokenHash: hashedRefreshToken,
	}

	orgAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}
	hasOrganization := len(orgAccesses) > 0

	var token, xsrfToken string
	var expiration time.Duration
	// Create the user in a transaction
	err = s.Repository.RunTransaction(func(tx port.Transaction) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		expiration = auth.TmpTokenExpiration
		if !config.Config.IsCloudEdition {
			// For self-hosted, create initial organization
			if err := s.createInitialOrganizationForSelfHosted(ctx, tx, u); err != nil {
				return err
			}
			expiration = auth.TokenExpiration()
			hasOrganization = true
		}

		// Create token
		var err error
		token, xsrfToken, _, _, _, err = s.createTokens(u.ID, expiration)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.RegisterWithMagicLinkOutput{
		Token:           token,
		RefreshToken:    plainRefreshToken,
		XSRFToken:       xsrfToken,
		ExpiresAt:       strconv.FormatInt(now.Add(expiration).Unix(), 10),
		HasOrganization: hasOrganization,
	}, nil
}

// RequestInvitationMagicLink sends a magic link for invitation authentication.
func (s *ServiceCE) RequestInvitationMagicLink(ctx context.Context, in dto.RequestInvitationMagicLinkInput) (*dto.RequestInvitationMagicLinkOutput, error) {
	// Parse and validate invitation token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.InvitationToken)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get invitation
	userInvitation, err := s.Repository.User().GetInvitation(ctx, user.InvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get organization
	invitedOrg, err := s.Repository.Organization().Get(ctx, organization.ByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	// Verify organization access in cloud edition
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		hostOrg, err := s.Repository.Organization().Get(ctx, organization.BySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	// Check if user exists
	exists, err := s.Repository.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	// Create magic link token
	tok, err := createInvitationMagicLinkToken(c.Email)
	if err != nil {
		return nil, err
	}

	// Build magic link URL
	url, err := buildInvitationMagicLinkURL(ptrconv.SafeValue(invitedOrg.Subdomain), tok)
	if err != nil {
		return nil, err
	}

	// Send magic link email
	if err := s.sendInvitationMagicLinkEmail(ctx, c.Email, "there", url); err != nil {
		return nil, err
	}

	return &dto.RequestInvitationMagicLinkOutput{
		Email: c.Email,
		IsNew: !exists,
	}, nil
}

// AuthenticateWithInvitationMagicLink authenticates a user with an invitation magic link.
func (s *ServiceCE) AuthenticateWithInvitationMagicLink(ctx context.Context, in dto.AuthenticateWithInvitationMagicLinkInput) (*dto.AuthenticateWithInvitationMagicLinkOutput, error) {
	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectInvitationMagicLink {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get invitation
	userInvitation, err := s.Repository.User().GetInvitation(ctx, user.InvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get organization
	invitedOrg, err := s.Repository.Organization().Get(ctx, organization.ByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		hostOrg, err := s.Repository.Organization().Get(ctx, organization.BySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = ptrconv.SafeValue(hostOrg.Subdomain)
	}

	// Check if user exists
	exists, err := s.Repository.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := createMagicLinkRegistrationToken(c.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to generate registration token: %w", err)
		}

		return &dto.AuthenticateWithInvitationMagicLinkOutput{
			Token:     registrationToken,
			IsNewUser: true,
		}, nil
	}

	// Get existing user
	u, err := s.Repository.User().Get(ctx, user.ByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Create organization access
	orgAccess := &user.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	// Generate token and refresh token
	now := time.Now()
	expiresAt := now.Add(auth.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err = s.Repository.RunTransaction(func(tx port.Transaction) error {
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
		return nil, err
	}

	return &dto.AuthenticateWithInvitationMagicLinkOutput{
		AuthURL:   config.Config.OrgBaseURL(orgSubdomain) + auth.SaveAuthPath,
		Token:     token,
		Domain:    config.Config.OrgDomain(orgSubdomain),
		IsNewUser: false,
	}, nil
}

// RegisterWithInvitationMagicLink registers a new user with an invitation magic link.
func (s *ServiceCE) RegisterWithInvitationMagicLink(ctx context.Context, in dto.RegisterWithInvitationMagicLinkInput) (*dto.RegisterWithInvitationMagicLinkOutput, error) {
	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectMagicLinkRegistration {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get invitation
	userInvitation, err := s.Repository.User().GetInvitation(ctx, user.InvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get organization
	invitedOrg, err := s.Repository.Organization().Get(ctx, organization.ByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		hostOrg, err := s.Repository.Organization().Get(ctx, organization.BySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = ptrconv.SafeValue(hostOrg.Subdomain)
	}

	// Generate refresh token
	plainRefreshToken, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	// Create new user
	now := time.Now()
	expiresAt := now.Add(auth.TokenExpiration())
	u := &user.User{
		ID:               uuid.Must(uuid.NewV4()),
		FirstName:        in.FirstName,
		LastName:         in.LastName,
		Email:            c.Email,
		RefreshTokenHash: hashedRefreshToken,
	}

	// Create organization access
	orgAccess := &user.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	// Generate token
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err = s.Repository.RunTransaction(func(tx port.Transaction) error {
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
		return nil, err
	}

	return &dto.RegisterWithInvitationMagicLinkOutput{
		Token:        token,
		RefreshToken: plainRefreshToken,
		XSRFToken:    xsrfToken,
		ExpiresAt:    strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:       config.Config.OrgDomain(orgSubdomain),
	}, nil
}

// RequestGoogleAuthLink sends a Google Auth link for invitation authentication.
func (s *ServiceCE) RequestGoogleAuthLink(ctx context.Context) (*dto.RequestGoogleAuthLinkOutput, error) {
	var hostSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
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
		return nil, err
	}

	googleOAuthClient := newGoogleOAuthClient()
	url, err := googleOAuthClient.getGoogleAuthCodeURL(ctx, stateToken)
	if err != nil {
		return nil, err
	}

	return &dto.RequestGoogleAuthLinkOutput{
		AuthURL: url,
	}, nil
}

func (s *ServiceCE) AuthenticateWithGoogle(ctx context.Context, in dto.AuthenticateWithGoogleInput) (*dto.AuthenticateWithGoogleOutput, error) {
	// Parse and validate state token
	stateClaims, err := jwt.ParseToken[*jwt.UserGoogleAuthLinkClaims](in.State)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	if stateClaims.Subject != jwt.UserSignatureSubjectGoogleAuthLink {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get Google token and user info
	googleOAuthClient := newGoogleOAuthClient()
	tok, err := googleOAuthClient.getGoogleToken(ctx, in.Code)
	if err != nil {
		return nil, err
	}

	userInfo, err := googleOAuthClient.getGoogleUserInfo(ctx, tok)
	if err != nil {
		return nil, err
	}

	// In staging environment, only allow @trysourcetool.com email addresses
	if config.Config.Env == config.EnvStaging && !strings.HasSuffix(userInfo.email, "@trysourcetool.com") {
		return nil, errdefs.ErrPermissionDenied(errors.New("access restricted in staging environment"))
	}

	// Check if user exists
	exists, err := s.Repository.User().IsEmailExists(ctx, userInfo.email)
	if err != nil {
		return nil, err
	}

	if !exists {
		if !config.Config.IsCloudEdition && stateClaims.Flow == jwt.GoogleAuthFlowStandard {
			if err := s.validateSelfHostedOrganization(ctx); err != nil {
				return nil, err
			}
		}

		var role string
		if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
			// Verify invitation exists
			userInvitation, err := s.Repository.User().GetInvitation(ctx, user.InvitationByEmail(userInfo.email), user.InvitationByOrganizationID(stateClaims.InvitationOrgID))
			if err != nil {
				return nil, errdefs.ErrInvalidArgument(errors.New("invalid invitation"))
			}
			role = userInvitation.Role.String()
		}

		// Generate registration token with flow info
		registrationToken, err := createGoogleRegistrationToken(
			userInfo.id,
			userInfo.email,
			userInfo.givenName,
			userInfo.familyName,
			stateClaims.Flow,
			stateClaims.InvitationOrgID,
			role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create registration token: %w", err)
		}

		return &dto.AuthenticateWithGoogleOutput{
			Token:           registrationToken,
			IsNewUser:       true,
			HasOrganization: stateClaims.Flow == jwt.GoogleAuthFlowInvitation,
			Flow:            string(stateClaims.Flow),
			FirstName:       userInfo.givenName,
			LastName:        userInfo.familyName,
		}, nil
	}

	// For existing users
	u, err := s.Repository.User().Get(ctx, user.ByEmail(userInfo.email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	needsGoogleIDUpdate := u.GoogleID == ""

	var org *organization.Organization
	var orgAccess *user.UserOrganizationAccess
	var orgSubdomain string

	if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
		// Handle invitation flow for existing users
		invitedOrg, err := s.Repository.Organization().Get(ctx, organization.ByID(stateClaims.InvitationOrgID))
		if err != nil {
			return nil, fmt.Errorf("failed to get invited organization: %w", err)
		}

		userInvitation, err := s.Repository.User().GetInvitation(ctx,
			user.InvitationByEmail(userInfo.email),
			user.InvitationByOrganizationID(stateClaims.InvitationOrgID))
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(errors.New("invalid invitation"))
		}

		orgAccess = &user.UserOrganizationAccess{
			ID:             uuid.Must(uuid.NewV4()),
			UserID:         u.ID,
			OrganizationID: invitedOrg.ID,
			Role:           userInvitation.Role,
		}
		org = invitedOrg
		orgSubdomain = ptrconv.SafeValue(invitedOrg.Subdomain)
	} else {
		// Standard flow - get user's organization info
		// Get all organization accesses for the user
		orgAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByUserID(u.ID))
		if err != nil {
			return nil, err
		}

		if config.Config.IsCloudEdition {
			if len(orgAccesses) > 1 {
				hostSubdomain := stateClaims.HostSubdomain
				if hostSubdomain == "" {
					// Handle multiple organizations by sending email with login URLs
					loginURLs := make([]string, 0, len(orgAccesses))
					for _, access := range orgAccesses {
						org, err := s.Repository.Organization().Get(ctx, organization.ByID(access.OrganizationID))
						if err != nil {
							return nil, err
						}

						url, err := buildLoginURL(ptrconv.SafeValue(org.Subdomain))
						if err != nil {
							return nil, err
						}
						loginURLs = append(loginURLs, url)
					}

					// Send email with multiple organization links
					if err := s.sendMultipleOrganizationsLoginEmail(ctx, u.Email, u.FirstName, loginURLs); err != nil {
						return nil, err
					}

					return &dto.AuthenticateWithGoogleOutput{
						IsNewUser:                false,
						HasOrganization:          true,
						HasMultipleOrganizations: true,
						Flow:                     string(stateClaims.Flow),
					}, nil
				} else {
					org, err = s.Repository.Organization().Get(ctx, organization.BySubdomain(hostSubdomain))
					if err != nil {
						return nil, err
					}
					orgAccess, err = s.Repository.User().GetOrganizationAccess(ctx,
						user.OrganizationAccessByUserID(u.ID),
						user.OrganizationAccessByOrganizationID(org.ID),
					)
					if err != nil {
						return nil, err
					}
					orgSubdomain = ptrconv.SafeValue(org.Subdomain)
				}
			} else {
				// Single organization case
				orgAccess = orgAccesses[0]

				org, err = s.Repository.Organization().Get(ctx, organization.ByID(orgAccess.OrganizationID))
				if err != nil {
					return nil, err
				}
				orgSubdomain = ptrconv.SafeValue(org.Subdomain)
			}
		} else {
			// Self-hosted mode
			orgAccess = orgAccesses[0]
			org, err = s.Repository.Organization().Get(ctx, organization.ByID(orgAccess.OrganizationID))
			if err != nil {
				return nil, err
			}
		}
	}

	// Generate temporary auth tokens
	token, xsrfToken, plainRefreshToken, hashedRefreshToken, _, err := s.createTokens(u.ID, auth.TmpTokenExpiration)
	if err != nil {
		return nil, err
	}

	u.RefreshTokenHash = hashedRefreshToken
	if needsGoogleIDUpdate {
		u.GoogleID = userInfo.id
	}

	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return nil, err
	}

	if err = s.Repository.RunTransaction(func(tx port.Transaction) error {
		if stateClaims.Flow == jwt.GoogleAuthFlowInvitation {
			// For invitation flow, create org access and delete invitation
			userInvitation, err := s.Repository.User().GetInvitation(ctx,
				user.InvitationByEmail(userInfo.email),
				user.InvitationByOrganizationID(stateClaims.InvitationOrgID))
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
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.AuthenticateWithGoogleOutput{
		AuthURL:         authURL,
		Token:           token,
		HasOrganization: orgAccess != nil,
		RefreshToken:    plainRefreshToken,
		XSRFToken:       xsrfToken,
		Domain:          config.Config.OrgDomain(orgSubdomain),
		IsNewUser:       false,
		Flow:            string(stateClaims.Flow),
	}, nil
}

// RegisterWithGoogle registers a new user based on the token received after Google OAuth confirmation.
func (s *ServiceCE) RegisterWithGoogle(ctx context.Context, in dto.RegisterWithGoogleInput) (*dto.RegisterWithGoogleOutput, error) {
	// Parse and validate registration token
	claims, err := jwt.ParseToken[*jwt.UserGoogleRegistrationClaims](in.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid registration token: %w", err)
	}
	if claims.Subject != jwt.UserSignatureSubjectGoogleRegistration {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject for google registration"))
	}

	// Check if user already exists
	exists, err := s.Repository.User().IsEmailExists(ctx, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errdefs.ErrUserEmailAlreadyExists(fmt.Errorf("user with email %s already exists", claims.Email))
	}

	plainRefreshToken, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	now := time.Now()
	tokenExpiration := auth.TokenExpiration()
	expiresAt := now.Add(tokenExpiration)
	u := &user.User{
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
	err = s.Repository.RunTransaction(func(tx port.Transaction) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if claims.Flow == jwt.GoogleAuthFlowInvitation {
			invitedOrg, err := s.Repository.Organization().Get(ctx, organization.ByID(claims.InvitationOrgID))
			if err != nil {
				return fmt.Errorf("failed to get invited organization: %w", err)
			}

			userInvitation, err := s.Repository.User().GetInvitation(ctx,
				user.InvitationByEmail(claims.Email),
				user.InvitationByOrganizationID(claims.InvitationOrgID))
			if err != nil {
				return fmt.Errorf("failed to get invitation: %w", err)
			}

			orgAccess := &user.UserOrganizationAccess{
				ID:             uuid.Must(uuid.NewV4()),
				UserID:         u.ID,
				OrganizationID: claims.InvitationOrgID,
				Role:           user.UserOrganizationRoleFromString(claims.Role),
			}

			if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
				return fmt.Errorf("failed to delete invitation: %w", err)
			}

			if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
				return fmt.Errorf("failed to create organization access: %w", err)
			}

			if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
				return fmt.Errorf("failed to create personal API key: %w", err)
			}

			orgSubdomain = ptrconv.SafeValue(invitedOrg.Subdomain)
			hasOrganization = true
		} else {
			if !config.Config.IsCloudEdition {
				if err := s.createInitialOrganizationForSelfHosted(ctx, tx, u); err != nil {
					return fmt.Errorf("failed to create initial organization: %w", err)
				}
				hasOrganization = true
			}
		}

		var err error
		token, xsrfToken, _, _, _, err = s.createTokens(u.ID, tokenExpiration)
		if err != nil {
			return fmt.Errorf("failed to create auth token: %w", err)
		}

		if hasOrganization {
			authURL, err = buildSaveAuthURL(orgSubdomain)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.RegisterWithGoogleOutput{
		Token:           token,
		RefreshToken:    plainRefreshToken,
		XSRFToken:       xsrfToken,
		ExpiresAt:       strconv.FormatInt(expiresAt.Unix(), 10),
		AuthURL:         authURL,
		HasOrganization: hasOrganization,
	}, nil
}

// RequestInvitationGoogleAuthLink prepares the Google OAuth URL for an invited user.
func (s *ServiceCE) RequestInvitationGoogleAuthLink(ctx context.Context, in dto.RequestInvitationGoogleAuthLinkInput) (*dto.RequestInvitationGoogleAuthLinkOutput, error) {
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.InvitationToken)
	if err != nil {
		return nil, fmt.Errorf("invalid invitation token: %w", err)
	}
	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject for invitation"))
	}

	userInvitation, err := s.Repository.User().GetInvitation(ctx, user.InvitationByEmail(c.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve invitation: %w", err)
	}

	invitedOrg, err := s.Repository.Organization().Get(ctx, organization.ByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve invited organization: %w", err)
	}

	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		if subdomain == "" || subdomain == "auth" {
			return nil, errdefs.ErrInvalidArgument(errors.New("invitation must be accessed via organization subdomain"))
		}
		hostOrg, err := s.Repository.Organization().Get(ctx, organization.BySubdomain(subdomain))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve host organization: %w", err)
		}
		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invitation organization mismatch"))
		}
	}

	var hostSubdomain string
	if config.Config.IsCloudEdition {
		hostSubdomain = ctxutil.Subdomain(ctx)
	}

	stateToken, err := createGoogleAuthLinkToken(
		jwt.GoogleAuthFlowInvitation,
		invitedOrg.ID,
		hostSubdomain,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create state token: %w", err)
	}

	googleOAuthClient := newGoogleOAuthClient()
	url, err := googleOAuthClient.getGoogleAuthCodeURL(ctx, stateToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get google auth code url: %w", err)
	}

	return &dto.RequestInvitationGoogleAuthLinkOutput{
		AuthURL: url,
	}, nil
}

func (s *ServiceCE) RefreshToken(ctx context.Context, in dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error) {
	// Validate XSRF token consistency
	if in.XSRFTokenCookie != in.XSRFTokenHeader {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid xsrf token"))
	}

	// Get user by refresh token
	hashedRefreshToken := hashRefreshToken(in.RefreshToken)
	u, err := s.Repository.User().Get(ctx, user.ByRefreshTokenHash(hashedRefreshToken))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	// Get current subdomain and resolve organization
	subdomain := ctxutil.Subdomain(ctx)
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// Verify user has access to this organization
			_, _, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return nil, err
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
	expiresAt := now.Add(auth.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	return &dto.RefreshTokenOutput{
		Token:        token,
		RefreshToken: in.RefreshToken,
		XSRFToken:    xsrfToken,
		ExpiresAt:    strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:       config.Config.OrgDomain(orgSubdomain),
	}, nil
}

func (s *ServiceCE) Logout(ctx context.Context) (*dto.LogoutOutput, error) {
	u := ctxutil.CurrentUser(ctx)

	orgAccessOpts := []user.OrganizationAccessQuery{
		user.OrganizationAccessByUserID(u.ID),
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		orgAccessOpts = append(orgAccessOpts, user.OrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := s.Repository.User().GetOrganizationAccess(ctx, orgAccessOpts...)
	if err != nil {
		return nil, err
	}

	return &dto.LogoutOutput{
		Domain: config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) Save(ctx context.Context, in dto.SaveAuthInput) (*dto.SaveAuthOutput, error) {
	// Parse and validate token
	c, err := jwt.ParseToken[*jwt.UserAuthClaims](in.Token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	u, err := s.Repository.User().Get(ctx, user.ByID(userID))
	if err != nil {
		return nil, err
	}

	// Get current subdomain and verify organization access
	subdomain := ctxutil.Subdomain(ctx)
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain, verify user has access
			_, _, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return nil, err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain, use default
			orgSubdomain = "auth"
		}
	}

	// Generate token and refresh token
	now := time.Now()
	expiresAt := now.Add(auth.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	plainRefreshToken, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	// Update user's refresh token
	u.RefreshTokenHash = hashedRefreshToken

	// Save changes
	if err = s.Repository.RunTransaction(func(tx port.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.SaveAuthOutput{
		Token:        token,
		RefreshToken: plainRefreshToken,
		XSRFToken:    xsrfToken,
		ExpiresAt:    strconv.FormatInt(expiresAt.Unix(), 10),
		RedirectURL:  config.Config.OrgBaseURL(orgSubdomain),
		Domain:       config.Config.OrgDomain(orgSubdomain),
	}, nil
}

func (s *ServiceCE) ObtainAuthToken(ctx context.Context) (*dto.ObtainAuthTokenOutput, error) {
	// Get current user from context
	u := ctxutil.CurrentUser(ctx)
	if u == nil {
		return nil, errdefs.ErrUnauthenticated(errors.New("no user in context"))
	}

	// Get user's organization info
	org, _, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Generate temporary token
	now := time.Now()
	expiresAt := now.Add(auth.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	// Build auth URL with organization subdomain
	authURL, err := buildSaveAuthURL(ptrconv.SafeValue(org.Subdomain))
	if err != nil {
		return nil, err
	}

	// Update user
	if err = s.Repository.RunTransaction(func(tx port.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.ObtainAuthTokenOutput{
		AuthURL: authURL,
		Token:   token,
	}, nil
}

func (s *ServiceCE) createPersonalAPIKey(ctx context.Context, tx port.Transaction, u *user.User, org *organization.Organization) error {
	devEnv, err := s.Repository.Environment().Get(ctx, environment.ByOrganizationID(org.ID), environment.BySlug(environment.EnvironmentSlugDevelopment))
	if err != nil {
		return err
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}

	apiKey := &apikey.APIKey{
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
func (s *ServiceCE) createTokens(userID uuid.UUID, expiration time.Duration) (token, xsrfToken, plainRefreshToken, hashedRefreshToken string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(expiration)
	xsrfToken = uuid.Must(uuid.NewV4()).String()

	token, err = createAuthToken(userID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return "", "", "", "", time.Time{}, err
	}

	plainRefreshToken, hashedRefreshToken, err = generateRefreshToken()
	if err != nil {
		return "", "", "", "", time.Time{}, errdefs.ErrInternal(err)
	}

	return token, xsrfToken, plainRefreshToken, hashedRefreshToken, expiresAt, nil
}

// resolveOrganizationBySubdomain gets an organization by subdomain and verifies the user has access.
// Deprecated: Use getOrganizationBySubdomain instead.
func (s *ServiceCE) resolveOrganizationBySubdomain(ctx context.Context, u *user.User, subdomain string) (*organization.Organization, *user.UserOrganizationAccess, error) {
	if subdomain == "" {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("subdomain cannot be empty"))
	}

	return s.getOrganizationBySubdomain(ctx, u, subdomain)
}

// validateSelfHostedOrganization checks if creating a new organization is allowed in self-hosted mode.
func (s *ServiceCE) validateSelfHostedOrganization(ctx context.Context) error {
	if !config.Config.IsCloudEdition {
		// In self-hosted mode, check if an organization already exists
		if _, err := s.Repository.Organization().Get(ctx); err == nil {
			return errdefs.ErrPermissionDenied(errors.New("only one organization is allowed in self-hosted edition"))
		}
	}
	return nil
}

func (s *ServiceCE) createInitialOrganizationForSelfHosted(ctx context.Context, tx port.Transaction, u *user.User) error {
	if config.Config.IsCloudEdition {
		return nil
	}

	org := &organization.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: nil, // Empty subdomain for non-cloud edition
	}
	if err := tx.Organization().Create(ctx, org); err != nil {
		return err
	}

	orgAccess := &user.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: org.ID,
		Role:           user.UserOrganizationRoleAdmin,
	}
	if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
		return err
	}

	devEnv := &environment.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		Name:           environment.EnvironmentNameDevelopment,
		Slug:           environment.EnvironmentSlugDevelopment,
		Color:          environment.EnvironmentColorDevelopment,
	}
	envs := []*environment.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: org.ID,
			Name:           environment.EnvironmentNameProduction,
			Slug:           environment.EnvironmentSlugProduction,
			Color:          environment.EnvironmentColorProduction,
		},
		devEnv,
	}
	if err := tx.Environment().BulkInsert(ctx, envs); err != nil {
		return err
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}
	apiKey := &apikey.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         u.ID,
		Name:           "",
		Key:            key,
	}
	if err := tx.APIKey().Create(ctx, apiKey); err != nil {
		return err
	}

	return nil
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *ServiceCE) getUserOrganizationInfo(ctx context.Context) (*organization.Organization, *user.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, ctxutil.CurrentUser(ctx))
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

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *ServiceCE) getOrganizationInfo(ctx context.Context, u *user.User) (*organization.Organization, *user.UserOrganizationAccess, error) {
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
