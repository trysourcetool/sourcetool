package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/jwt"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

// Service defines the interface for user-related operations.
type Service interface {
	// User management methods
	GetMe(context.Context) (*dto.GetMeOutput, error)
	List(context.Context) (*dto.ListUsersOutput, error)
	Update(context.Context, dto.UpdateUserInput) (*dto.UpdateUserOutput, error)

	// Email operations
	SendUpdateEmailInstructions(context.Context, dto.SendUpdateUserEmailInstructionsInput) error
	UpdateEmail(context.Context, dto.UpdateUserEmailInput) (*dto.UpdateUserEmailOutput, error)

	// Passwordless Authentication methods
	RequestMagicLink(context.Context, dto.RequestMagicLinkInput) (*dto.RequestMagicLinkOutput, error)
	AuthenticateWithMagicLink(context.Context, dto.AuthenticateWithMagicLinkInput) (*dto.AuthenticateWithMagicLinkOutput, error)
	RegisterWithMagicLink(context.Context, dto.RegisterWithMagicLinkInput) (*dto.RegisterWithMagicLinkOutput, error)

	// Invitation Magic Link methods
	RequestInvitationMagicLink(context.Context, dto.RequestInvitationMagicLinkInput) (*dto.RequestInvitationMagicLinkOutput, error)
	AuthenticateWithInvitationMagicLink(context.Context, dto.AuthenticateWithInvitationMagicLinkInput) (*dto.AuthenticateWithInvitationMagicLinkOutput, error)
	RegisterWithInvitationMagicLink(context.Context, dto.RegisterWithInvitationMagicLinkInput) (*dto.RegisterWithInvitationMagicLinkOutput, error)

	// Authentication methods
	SignOut(context.Context) (*dto.SignOutOutput, error)
	RefreshToken(context.Context, dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error)
	SaveAuth(context.Context, dto.SaveAuthInput) (*dto.SaveAuthOutput, error)
	ObtainAuthToken(context.Context) (*dto.ObtainAuthTokenOutput, error)

	// Invitation methods
	Invite(context.Context, dto.InviteUsersInput) (*dto.InviteUsersOutput, error)
	ResendInvitation(context.Context, dto.ResendInvitationInput) (*dto.ResendInvitationOutput, error)
	SignInInvitation(context.Context, dto.SignInInvitationInput) (*dto.SignInInvitationOutput, error)
	SignUpInvitation(context.Context, dto.SignUpInvitationInput) (*dto.SignUpInvitationOutput, error)

	// Google OAuth methods
	GetGoogleAuthCodeURL(context.Context) (*dto.GetGoogleAuthCodeURLOutput, error)
	GoogleOAuthCallback(context.Context, dto.GoogleOAuthCallbackInput) (*dto.GoogleOAuthCallbackOutput, error)
	SignInWithGoogle(context.Context, dto.SignInWithGoogleInput) (*dto.SignInWithGoogleOutput, error)
	SignUpWithGoogle(context.Context, dto.SignUpWithGoogleInput) (*dto.SignUpWithGoogleOutput, error)
	GetGoogleAuthCodeURLInvitation(context.Context, dto.GetGoogleAuthCodeURLInvitationInput) (*dto.GetGoogleAuthCodeURLInvitationOutput, error)
	SignInWithGoogleInvitation(context.Context, dto.SignInWithGoogleInvitationInput) (*dto.SignInWithGoogleInvitationOutput, error)
	SignUpWithGoogleInvitation(context.Context, dto.SignUpWithGoogleInvitationInput) (*dto.SignUpWithGoogleInvitationOutput, error)
}

// ServiceCE implements the Service interface for the Community Edition.
type ServiceCE struct {
	*infra.Dependency
}

// NewServiceCE creates a new instance of the ServiceCE.
func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) GetMe(ctx context.Context) (*dto.GetMeOutput, error) {
	currentUser := ctxutil.CurrentUser(ctx)
	currentOrg := ctxutil.CurrentOrganization(ctx)
	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByUserID(currentUser.ID))
	if err != nil {
		return nil, err
	}
	var orgAccess *model.UserOrganizationAccess
	if len(orgAccesses) > 0 {
		if currentOrg == nil {
			if len(orgAccesses) > 1 {
				return nil, errdefs.ErrUserMultipleOrganizations(errors.New("user has multiple organizations"))
			}
			currentOrg, err = s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccesses[0].OrganizationID))
			if err != nil {
				return nil, err
			}
			orgAccess = orgAccesses[0]
		} else {
			var err error
			orgAccess, err = s.Store.User().GetOrganizationAccess(ctx,
				storeopts.UserOrganizationAccessByUserID(currentUser.ID),
				storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID))
			if err != nil {
				return nil, err
			}
		}
	}

	var role model.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.GetMeOutput{
		User: dto.UserFromModel(currentUser, currentOrg, role),
	}, nil
}

func (s *ServiceCE) List(ctx context.Context) (*dto.ListUsersOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)

	users, err := s.Store.User().List(ctx, storeopts.UserByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userInvitations, err := s.Store.User().ListInvitations(ctx, storeopts.UserInvitationByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}
	roleMap := make(map[uuid.UUID]model.UserOrganizationRole)
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

	var role model.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.UpdateUserOutput{
		User: dto.UserFromModel(currentUser, org, role),
	}, nil
}

// SendUpdateEmailInstructions sends instructions for updating a user's email address.
func (s *ServiceCE) SendUpdateEmailInstructions(ctx context.Context, in dto.SendUpdateUserEmailInstructionsInput) error {
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
	tok, err := createUserToken(currentUser.ID.String(), in.Email, time.Now().Add(model.EmailTokenExpiration), jwt.UserSignatureSubjectUpdateEmail)
	if err != nil {
		return err
	}

	// Build update URL
	url, err := buildUpdateEmailURL(conv.SafeValue(currentOrg.Subdomain), tok)
	if err != nil {
		return err
	}

	return s.Mailer.User().SendUpdateEmailInstructions(ctx, &model.SendUpdateUserEmailInstructions{
		To:        in.Email,
		FirstName: currentUser.FirstName,
		URL:       url,
	})
}

func (s *ServiceCE) UpdateEmail(ctx context.Context, in dto.UpdateUserEmailInput) (*dto.UpdateUserEmailOutput, error) {
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
	u, err := s.Store.User().Get(ctx, storeopts.UserByID(userID))
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

	var role model.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.UpdateUserEmailOutput{
		User: dto.UserFromModel(currentUser, org, role),
	}, nil
}

// Works for both existing users (login) and new users (signup).
func (s *ServiceCE) RequestMagicLink(ctx context.Context, in dto.RequestMagicLinkInput) (*dto.RequestMagicLinkOutput, error) {
	// Check if email exists
	exists, err := s.Store.User().IsEmailExists(ctx, in.Email)
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
			org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
			if err != nil {
				return nil, err
			}

			if exists {
				// For existing users, check if they have access to this organization
				u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(in.Email))
				if err != nil {
					return nil, err
				}

				_, err = s.Store.User().GetOrganizationAccess(ctx,
					storeopts.UserOrganizationAccessByUserID(u.ID),
					storeopts.UserOrganizationAccessByOrganizationID(org.ID))
				if err != nil {
					return nil, errdefs.ErrUnauthenticated(errors.New("user does not have access to this organization"))
				}
				firstName = u.FirstName
			} else {
				// For new users, registration is only allowed through invitations
				return nil, errdefs.ErrPermissionDenied(errors.New("registration is only allowed through invitations"))
			}
		}
	}

	if exists && !config.Config.IsCloudEdition {
		// Get user by email for existing users
		u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(in.Email))
		if err != nil {
			return nil, err
		}
		firstName = u.FirstName

		// Get user's organization access information
		orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByUserID(u.ID))
		if err != nil {
			return nil, err
		}

		// Cloud edition specific handling for multiple organizations
		if config.Config.IsCloudEdition && len(orgAccesses) > 1 {
			// Handle multiple organizations
			loginURLs := make([]string, 0, len(orgAccesses))
			for _, access := range orgAccesses {
				org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(access.OrganizationID))
				if err != nil {
					return nil, err
				}

				// Create org-specific magic link
				tok, err := createUserEmailToken(in.Email, time.Now().Add(15*time.Minute), jwt.UserSignatureSubjectMagicLink)
				if err != nil {
					return nil, err
				}

				url, err := buildMagicLinkURL(conv.SafeValue(org.Subdomain), tok)
				if err != nil {
					return nil, err
				}
				loginURLs = append(loginURLs, url)
			}

			// Send email with multiple organization links
			if err := s.Mailer.User().SendMultipleOrganizationsMagicLinkEmail(ctx, &model.SendMultipleOrganizationsMagicLinkEmail{
				To:        in.Email,
				FirstName: firstName,
				Email:     in.Email,
				LoginURLs: loginURLs,
			}); err != nil {
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
	tok, err := createUserEmailToken(in.Email, time.Now().Add(15*time.Minute), jwt.UserSignatureSubjectMagicLink)
	if err != nil {
		return nil, err
	}

	// Build magic link URL
	url, err := buildMagicLinkURL(subdomain, tok)
	if err != nil {
		return nil, err
	}

	// Send magic link email
	if err := s.Mailer.User().SendMagicLinkEmail(ctx, &model.SendMagicLinkEmail{
		To:        in.Email,
		FirstName: firstName,
		URL:       url,
	}); err != nil {
		return nil, err
	}

	return &dto.RequestMagicLinkOutput{
		Email: in.Email,
		IsNew: isNewUser,
	}, nil
}

// Can handle both login and signup based on input parameters.
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
	exists, err := s.Store.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := createRegistrationToken(c.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to generate registration token: %w", err)
		}

		return &dto.AuthenticateWithMagicLinkOutput{
			Token:                registrationToken,
			IsNewUser:            true,
			IsOrganizationExists: false,
		}, nil
	}

	// Get existing user
	user, err := s.Store.User().Get(ctx, storeopts.UserByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get user's organization access information
	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByUserID(user.ID))
	if err != nil {
		return nil, err
	}

	// Handle organization subdomain logic
	subdomain := ctxutil.Subdomain(ctx)
	var orgAccess *model.UserOrganizationAccess
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain, resolve org and access
			_, orgAccess, err = s.resolveOrganizationBySubdomain(ctx, user, subdomain)
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
				org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
				if err != nil {
					return nil, err
				}
				orgSubdomain = conv.SafeValue(org.Subdomain)
			} else {
				return nil, errdefs.ErrUserMultipleOrganizations(errors.New("user has multiple organizations"))
			}
		}
	} else {
		// Self-hosted mode has only one organization
		orgAccess = orgAccesses[0]
		_, err = s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
		if err != nil {
			return nil, err
		}
	}

	// Create token, secret, etc.
	token, xsrfToken, plainSecret, hashedSecret, _, err := s.createTokenWithSecret(
		user.ID, model.TmpTokenExpiration)
	if err != nil {
		return nil, err
	}

	// Update user with new secret
	user.Secret = hashedSecret
	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, user)
	}); err != nil {
		return nil, err
	}

	return &dto.AuthenticateWithMagicLinkOutput{
		AuthURL:              authURL,
		Token:                token,
		IsOrganizationExists: orgAccess != nil,
		Secret:               plainSecret,
		XSRFToken:            xsrfToken,
		Domain:               config.Config.OrgDomain(orgSubdomain),
		IsNewUser:            false,
	}, nil
}

// RegisterWithMagicLink registers a new user with a magic link.
func (s *ServiceCE) RegisterWithMagicLink(ctx context.Context, in dto.RegisterWithMagicLinkInput) (*dto.RegisterWithMagicLinkOutput, error) {
	// Parse and validate the registration token
	claims, err := jwt.ParseToken[*jwt.UserMagicLinkRegistrationClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if claims.Subject != jwt.UserSignatureSubjectMagicLinkRegistration {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Generate secret and XSRF token
	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, err
	}

	// Create a new user
	user := &model.User{
		ID:        uuid.Must(uuid.NewV4()),
		Email:     claims.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Secret:    hashedSecret,
	}

	var token, xsrfToken string
	// Create the user in a transaction
	err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().Create(ctx, user); err != nil {
			return err
		}

		expiration := model.TmpTokenExpiration
		if !config.Config.IsCloudEdition {
			// For self-hosted, create initial organization
			if err := s.createInitialOrganizationForSelfHosted(ctx, tx, user); err != nil {
				return err
			}
			expiration = model.TokenExpiration()
		}

		// Create token
		var err error
		token, xsrfToken, _, _, _, err = s.createTokenWithSecret(user.ID, expiration)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.RegisterWithMagicLinkOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
	}, nil
}

// resolveOrganizationBySubdomain gets an organization by subdomain and verifies the user has access.
// Deprecated: Use getOrganizationBySubdomain instead.
func (s *ServiceCE) resolveOrganizationBySubdomain(ctx context.Context, u *model.User, subdomain string) (*model.Organization, *model.UserOrganizationAccess, error) {
	if subdomain == "" {
		return nil, nil, errdefs.ErrInvalidArgument(errors.New("subdomain cannot be empty"))
	}

	return s.getOrganizationBySubdomain(ctx, u, subdomain)
}

func (s *ServiceCE) SignInWithGoogle(ctx context.Context, in dto.SignInWithGoogleInput) (*dto.SignInWithGoogleOutput, error) {
	// Parse and validate the Google auth request token
	googleAuthReqClaims, err := jwt.ParseToken[*jwt.UserGoogleAuthRequestClaims](in.SessionToken)
	if err != nil {
		return nil, err
	}

	googleAuthReqID, err := uuid.FromString(googleAuthReqClaims.GoogleAuthRequestID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	// Get and validate the Google auth request
	googleAuthReq, err := s.Store.User().GetGoogleAuthRequest(ctx, googleAuthReqID)
	if err != nil {
		return nil, err
	}

	if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	// Get the user by email
	u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	// Get current subdomain and resolve organization access
	subdomain := ctxutil.Subdomain(ctx)
	var org *model.Organization
	var orgAccess *model.UserOrganizationAccess
	var orgSubdomain string

	// Handle organization lookup based on edition and subdomain
	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain, verify user has access
			_, orgAccess, err = s.resolveOrganizationBySubdomain(ctx, u, subdomain)
			if err != nil {
				return nil, err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain, get user's first organization access
			orgAccess, err = s.Store.User().GetOrganizationAccess(ctx,
				storeopts.UserOrganizationAccessByUserID(u.ID),
				storeopts.UserOrganizationAccessOrderBy("created_at DESC"))

			if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
				return nil, err
			}

			if orgAccess != nil {
				org, err = s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
				if err != nil {
					return nil, err
				}
				orgSubdomain = conv.SafeValue(org.Subdomain)
			} else {
				orgSubdomain = "auth"
			}
		}
	} else {
		// Self-hosted edition - get the user's organization
		_, orgAccess, err = s.getOrganizationInfo(ctx, u)
		if err != nil {
			if !errdefs.IsUserOrganizationAccessNotFound(err) {
				return nil, err
			}
			// If no organization access found, continue with nil orgAccess
		}
	}

	// Generate token and secret
	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	// Update user's secret
	u.Secret = hashedSecret

	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.SignInWithGoogleOutput{
		AuthURL:              authURL,
		Token:                token,
		Secret:               plainSecret,
		XSRFToken:            xsrfToken,
		IsOrganizationExists: orgAccess != nil,
		Domain:               config.Config.OrgDomain(orgSubdomain),
	}, nil
}

// validateSelfHostedOrganization checks if creating a new organization is allowed in self-hosted mode.
func (s *ServiceCE) validateSelfHostedOrganization(ctx context.Context) error {
	if !config.Config.IsCloudEdition {
		// In self-hosted mode, check if an organization already exists
		if _, err := s.Store.Organization().Get(ctx); err == nil {
			return errdefs.ErrPermissionDenied(errors.New("only one organization is allowed in self-hosted edition"))
		}
	}
	return nil
}

func (s *ServiceCE) createInitialOrganizationForSelfHosted(ctx context.Context, tx infra.Transaction, u *model.User) error {
	if config.Config.IsCloudEdition {
		return nil
	}

	org := &model.Organization{
		ID:        uuid.Must(uuid.NewV4()),
		Subdomain: nil, // Empty subdomain for non-cloud edition
	}
	if err := tx.Organization().Create(ctx, org); err != nil {
		return err
	}

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: org.ID,
		Role:           model.UserOrganizationRoleAdmin,
	}
	if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
		return err
	}

	devEnv := &model.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		Name:           model.EnvironmentNameDevelopment,
		Slug:           model.EnvironmentSlugDevelopment,
		Color:          model.EnvironmentColorDevelopment,
	}
	envs := []*model.Environment{
		{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: org.ID,
			Name:           model.EnvironmentNameProduction,
			Slug:           model.EnvironmentSlugProduction,
			Color:          model.EnvironmentColorProduction,
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
	apiKey := &model.APIKey{
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

func (s *ServiceCE) SignUpWithGoogle(ctx context.Context, in dto.SignUpWithGoogleInput) (*dto.SignUpWithGoogleOutput, error) {
	googleAuthReqClaims, err := jwt.ParseToken[*jwt.UserGoogleAuthRequestClaims](in.SessionToken)
	if err != nil {
		return nil, err
	}

	googleAuthReqID, err := uuid.FromString(googleAuthReqClaims.GoogleAuthRequestID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	googleAuthReq, err := s.Store.User().GetGoogleAuthRequest(ctx, googleAuthReqID)
	if err != nil {
		return nil, err
	}

	if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	requestUser, err := s.Store.User().GetRegistrationRequest(ctx, storeopts.UserRegistrationRequestByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	now := time.Now()
	u := &model.User{
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Email:                googleAuthReq.Email,
		Secret:               hashedSecret,
		EmailAuthenticatedAt: &now,
		GoogleID:             googleAuthReq.GoogleID,
	}

	var token string
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		if !config.Config.IsCloudEdition {
			if err := s.createInitialOrganizationForSelfHosted(ctx, tx, u); err != nil {
				return err
			}
			expiresAt := now.Add(model.TokenExpiration())
			token, err = createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
		} else {
			expiresAt := now.Add(model.TmpTokenExpiration)
			token, err = createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
		}
		if err != nil {
			return err
		}

		return tx.User().DeleteRegistrationRequest(ctx, requestUser)
	}); err != nil {
		return nil, err
	}

	return &dto.SignUpWithGoogleOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
	}, nil
}

func (s *ServiceCE) RefreshToken(ctx context.Context, in dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error) {
	// Validate XSRF token consistency
	if in.XSRFTokenCookie != in.XSRFTokenHeader {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid xsrf token"))
	}

	// Get user by secret
	hashedSecret := hashSecret(in.Secret)
	u, err := s.Store.User().Get(ctx, storeopts.UserBySecret(hashedSecret))
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
	expiresAt := now.Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	return &dto.RefreshTokenOutput{
		Token:     token,
		Secret:    in.Secret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    config.Config.OrgDomain(orgSubdomain),
	}, nil
}

func (s *ServiceCE) SaveAuth(ctx context.Context, in dto.SaveAuthInput) (*dto.SaveAuthOutput, error) {
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
	u, err := s.Store.User().Get(ctx, storeopts.UserByID(userID))
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

	// Generate token and secret
	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	// Update user's secret
	u.Secret = hashedSecret

	// Save changes
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.SaveAuthOutput{
		Token:       token,
		Secret:      plainSecret,
		XSRFToken:   xsrfToken,
		ExpiresAt:   strconv.FormatInt(expiresAt.Unix(), 10),
		RedirectURL: config.Config.OrgBaseURL(orgSubdomain),
		Domain:      config.Config.OrgDomain(orgSubdomain),
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
	expiresAt := now.Add(model.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	// Build auth URL with organization subdomain
	authURL, err := buildSaveAuthURL(conv.SafeValue(org.Subdomain))
	if err != nil {
		return nil, err
	}

	// Update user
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.ObtainAuthTokenOutput{
		AuthURL: authURL,
		Token:   token,
	}, nil
}

func (s *ServiceCE) Invite(ctx context.Context, in dto.InviteUsersInput) (*dto.InviteUsersOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditUser); err != nil {
		return nil, err
	}

	o := ctxutil.CurrentOrganization(ctx)
	u := ctxutil.CurrentUser(ctx)

	invitations := make([]*model.UserInvitation, 0)
	emailInput := &model.SendInvitationEmail{
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

		tok, err := createUserEmailToken(email, time.Now().Add(model.EmailTokenExpiration), jwt.UserSignatureSubjectInvitation)
		if err != nil {
			return nil, err
		}

		url, err := buildInvitationURL(conv.SafeValue(o.Subdomain), tok, email)
		if err != nil {
			return nil, err
		}

		emailInput.URLs[email] = url

		invitations = append(invitations, &model.UserInvitation{
			ID:             uuid.Must(uuid.NewV4()),
			OrganizationID: o.ID,
			Email:          email,
			Role:           model.UserOrganizationRoleFromString(in.Role),
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

	usersInvitationsOut := make([]*dto.UserInvitation, 0, len(invitations))
	for _, ui := range invitations {
		usersInvitationsOut = append(usersInvitationsOut, dto.UserInvitationFromModel(ui))
	}

	return &dto.InviteUsersOutput{
		UserInvitations: usersInvitationsOut,
	}, nil
}

func (s *ServiceCE) SignInInvitation(ctx context.Context, in dto.SignInInvitationInput) (*dto.SignInInvitationOutput, error) {
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.InvitationToken)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	expiresAt := time.Now().Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = hashedSecret

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
			return err
		}

		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.SignInInvitationOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) SignUpInvitation(ctx context.Context, in dto.SignUpInvitationInput) (*dto.SignUpInvitationOutput, error) {
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.InvitationToken)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	// Hash password
	hashedPassword, err := hashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	u := &model.User{
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Email:                c.Email,
		Password:             hashedPassword,
		Secret:               hashedSecret,
		EmailAuthenticatedAt: &now,
	}

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	var token string
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
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

		token, err = createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.SignUpInvitationOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) GetGoogleAuthCodeURL(ctx context.Context) (*dto.GetGoogleAuthCodeURLOutput, error) {
	// Check self-hosted organization restriction
	if err := s.validateSelfHostedOrganization(ctx); err != nil {
		return nil, err
	}

	state := uuid.Must(uuid.NewV4())
	googleOAuthClient := newGoogleOAuthClient()
	url, err := googleOAuthClient.getGoogleAuthCodeURL(ctx, state.String())
	if err != nil {
		return nil, err
	}

	googleAuthReqs, err := s.Store.User().ListExpiredGoogleAuthRequests(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().BulkDeleteGoogleAuthRequests(ctx, googleAuthReqs); err != nil {
			return err
		}

		return tx.User().CreateGoogleAuthRequest(ctx, &model.UserGoogleAuthRequest{
			ID:        state,
			Domain:    config.Config.OrgHostname("auth"),
			ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour),
			Invited:   false,
		})
	}); err != nil {
		return nil, err
	}

	return &dto.GetGoogleAuthCodeURLOutput{
		URL: url,
	}, nil
}

func (s *ServiceCE) GoogleOAuthCallback(ctx context.Context, in dto.GoogleOAuthCallbackInput) (*dto.GoogleOAuthCallbackOutput, error) {
	state, err := uuid.FromString(in.State)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	googleAuthReq, err := s.Store.User().GetGoogleAuthRequest(ctx, state)
	if err != nil {
		return nil, err
	}

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
		return &dto.GoogleOAuthCallbackOutput{
			SessionToken: "",
			IsUserExists: false,
			Domain:       googleAuthReq.Domain,
		}, nil
	}

	isUserExists, err := s.Store.User().IsEmailExists(ctx, userInfo.email)
	if err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
			return nil, err
		}

		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	googleAuthReq.GoogleID = userInfo.id
	googleAuthReq.Email = userInfo.email

	sessionToken, err := createGoogleAuthRequestToken(googleAuthReq.ID.String(), googleAuthReq.ExpiresAt, jwt.UserSignatureSubjectGoogleAuthRequest)
	if err != nil {
		return nil, err
	}

	if isUserExists {
		u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(userInfo.email))
		if err != nil {
			return nil, err
		}

		hasGoogleID := u.GoogleID != ""

		if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
			if !hasGoogleID {
				u.GoogleID = googleAuthReq.GoogleID
				if err := tx.User().Update(ctx, u); err != nil {
					return err
				}
			}

			return tx.User().UpdateGoogleAuthRequest(ctx, googleAuthReq)
		}); err != nil {
			return nil, err
		}

		return &dto.GoogleOAuthCallbackOutput{
			SessionToken: sessionToken,
			IsUserExists: isUserExists,
			Domain:       googleAuthReq.Domain,
		}, nil
	}

	requestExists, err := s.Store.User().IsRegistrationRequestExists(ctx, userInfo.email)
	if err != nil {
		return nil, err
	}

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if !requestExists {
			if err := tx.User().CreateRegistrationRequest(ctx, &model.UserRegistrationRequest{
				ID:    uuid.Must(uuid.NewV4()),
				Email: userInfo.email,
			}); err != nil {
				return err
			}
		}

		return tx.User().UpdateGoogleAuthRequest(ctx, googleAuthReq)
	}); err != nil {
		return nil, err
	}

	return &dto.GoogleOAuthCallbackOutput{
		SessionToken: sessionToken,
		IsUserExists: isUserExists,
		FirstName:    userInfo.givenName,
		LastName:     userInfo.familyName,
		Domain:       googleAuthReq.Domain,
		Invited:      googleAuthReq.Invited,
	}, nil
}

func (s *ServiceCE) GetGoogleAuthCodeURLInvitation(ctx context.Context, in dto.GetGoogleAuthCodeURLInvitationInput) (*dto.GetGoogleAuthCodeURLInvitationOutput, error) {
	state := uuid.Must(uuid.NewV4())
	googleOAuthClient := newGoogleOAuthClient()
	url, err := googleOAuthClient.getGoogleAuthCodeURL(ctx, state.String())
	if err != nil {
		return nil, err
	}

	googleAuthReqs, err := s.Store.User().ListExpiredGoogleAuthRequests(ctx)
	if err != nil {
		return nil, err
	}

	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.InvitationToken)
	if err != nil {
		return nil, err
	}
	if c.Subject != jwt.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().BulkDeleteGoogleAuthRequests(ctx, googleAuthReqs); err != nil {
			return err
		}

		return tx.User().CreateGoogleAuthRequest(ctx, &model.UserGoogleAuthRequest{
			ID:        state,
			Domain:    config.Config.OrgDomain(subdomain),
			ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour),
			Invited:   true,
		})
	}); err != nil {
		return nil, err
	}

	return &dto.GetGoogleAuthCodeURLInvitationOutput{
		URL: url,
	}, nil
}

func (s *ServiceCE) SignInWithGoogleInvitation(ctx context.Context, in dto.SignInWithGoogleInvitationInput) (*dto.SignInWithGoogleInvitationOutput, error) {
	googleAuthReqClaims, err := jwt.ParseToken[*jwt.UserGoogleAuthRequestClaims](in.SessionToken)
	if err != nil {
		return nil, err
	}

	googleAuthReqID, err := uuid.FromString(googleAuthReqClaims.GoogleAuthRequestID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	googleAuthReq, err := s.Store.User().GetGoogleAuthRequest(ctx, googleAuthReqID)
	if err != nil {
		return nil, err
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	expiresAt := time.Now().Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = hashedSecret

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().DeleteInvitation(ctx, userInvitation); err != nil {
			return err
		}

		if err := tx.User().CreateOrganizationAccess(ctx, orgAccess); err != nil {
			return err
		}

		if err := s.createPersonalAPIKey(ctx, tx, u, invitedOrg); err != nil {
			return err
		}

		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &dto.SignInWithGoogleInvitationOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) SignUpWithGoogleInvitation(ctx context.Context, in dto.SignUpWithGoogleInvitationInput) (*dto.SignUpWithGoogleInvitationOutput, error) {
	googleAuthReqClaims, err := jwt.ParseToken[*jwt.UserGoogleAuthRequestClaims](in.SessionToken)
	if err != nil {
		return nil, err
	}

	googleAuthReqID, err := uuid.FromString(googleAuthReqClaims.GoogleAuthRequestID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	googleAuthReq, err := s.Store.User().GetGoogleAuthRequest(ctx, googleAuthReqID)
	if err != nil {
		return nil, err
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	requestUser, err := s.Store.User().GetRegistrationRequest(ctx, storeopts.UserRegistrationRequestByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	u := &model.User{
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Email:                googleAuthReq.Email,
		Secret:               hashedSecret,
		EmailAuthenticatedAt: &now,
	}

	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	var token string
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
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

		token, err = createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
		if err != nil {
			return err
		}

		return tx.User().DeleteRegistrationRequest(ctx, requestUser)
	}); err != nil {
		return nil, err
	}

	return &dto.SignUpWithGoogleInvitationOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) ResendInvitation(ctx context.Context, in dto.ResendInvitationInput) (*dto.ResendInvitationOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditUser); err != nil {
		return nil, err
	}

	invitationID, err := uuid.FromString(in.InvitationID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByID(invitationID))
	if err != nil {
		return nil, err
	}

	o := ctxutil.CurrentOrganization(ctx)
	if userInvitation.OrganizationID != o.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	u := ctxutil.CurrentUser(ctx)

	tok, err := createUserEmailToken(userInvitation.Email, time.Now().Add(model.EmailTokenExpiration), jwt.UserSignatureSubjectInvitation)
	if err != nil {
		return nil, err
	}

	url, err := buildInvitationURL(conv.SafeValue(o.Subdomain), tok, userInvitation.Email)
	if err != nil {
		return nil, err
	}

	emailInput := &model.SendInvitationEmail{
		Invitees: u.FullName(),
		URLs:     map[string]string{userInvitation.Email: url},
	}

	if err := s.Mailer.User().SendInvitationEmail(ctx, emailInput); err != nil {
		return nil, err
	}

	return &dto.ResendInvitationOutput{
		UserInvitation: dto.UserInvitationFromModel(userInvitation),
	}, nil
}

func (s *ServiceCE) SignOut(ctx context.Context) (*dto.SignOutOutput, error) {
	u := ctxutil.CurrentUser(ctx)

	orgAccessOpts := []storeopts.UserOrganizationAccessOption{
		storeopts.UserOrganizationAccessByUserID(u.ID),
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		orgAccessOpts = append(orgAccessOpts, storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := s.Store.User().GetOrganizationAccess(ctx, orgAccessOpts...)
	if err != nil {
		return nil, err
	}

	return &dto.SignOutOutput{
		Domain: config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) createPersonalAPIKey(ctx context.Context, tx infra.Transaction, u *model.User, org *model.Organization) error {
	devEnv, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByOrganizationID(org.ID), storeopts.EnvironmentBySlug(model.EnvironmentSlugDevelopment))
	if err != nil {
		return err
	}

	key, err := devEnv.GenerateAPIKey()
	if err != nil {
		return err
	}

	apiKey := &model.APIKey{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: org.ID,
		EnvironmentID:  devEnv.ID,
		UserID:         u.ID,
		Name:           "",
		Key:            key,
	}

	return tx.APIKey().Create(ctx, apiKey)
}

// Helper functions for common operations

// createTokenWithSecret creates a new authentication token and secret.
func (s *ServiceCE) createTokenWithSecret(userID uuid.UUID, expiration time.Duration) (token, xsrfToken, plainSecret, hashedSecret string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(expiration)
	xsrfToken = uuid.Must(uuid.NewV4()).String()

	token, err = createAuthToken(userID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return "", "", "", "", time.Time{}, err
	}

	plainSecret, hashedSecret, err = generateSecret()
	if err != nil {
		return "", "", "", "", time.Time{}, errdefs.ErrInternal(err)
	}

	return token, xsrfToken, plainSecret, hashedSecret, expiresAt, nil
}

// getUserOrganizationInfo is a convenience wrapper that retrieves organization
// and access information for the current user from the context.
func (s *ServiceCE) getUserOrganizationInfo(ctx context.Context) (*model.Organization, *model.UserOrganizationAccess, error) {
	return s.getOrganizationInfo(ctx, ctxutil.CurrentUser(ctx))
}

// getOrganizationInfo retrieves organization and access information for the specified user.
// It handles both cloud and self-hosted editions with appropriate subdomain logic.
func (s *ServiceCE) getOrganizationInfo(ctx context.Context, u *model.User) (*model.Organization, *model.UserOrganizationAccess, error) {
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
func (s *ServiceCE) getOrganizationBySubdomain(ctx context.Context, u *model.User, subdomain string) (*model.Organization, *model.UserOrganizationAccess, error) {
	// Get organization by subdomain
	org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	// Verify user has access to this organization
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByOrganizationID(org.ID),
		storeopts.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// (typically the most recently created one).
func (s *ServiceCE) getDefaultOrganizationForUser(ctx context.Context, u *model.User) (*model.Organization, *model.UserOrganizationAccess, error) {
	// Get user's organization access
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByUserID(u.ID),
		storeopts.UserOrganizationAccessOrderBy("created_at DESC"))
	if err != nil {
		return nil, nil, err
	}

	// Get the organization
	org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, nil, err
	}

	return org, orgAccess, nil
}

// RequestInvitationMagicLink sends a magic link for invitation authentication
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
	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get organization
	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	// Verify organization access in cloud edition
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}
	}

	// Check if user exists
	exists, err := s.Store.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	// Create magic link token
	tok, err := createUserEmailToken(c.Email, time.Now().Add(15*time.Minute), jwt.UserSignatureSubjectInvitationMagicLink)
	if err != nil {
		return nil, err
	}

	// Build magic link URL
	url, err := buildInvitationMagicLinkURL(conv.SafeValue(invitedOrg.Subdomain), tok)
	if err != nil {
		return nil, err
	}

	// Send magic link email
	if err := s.Mailer.User().SendInvitationMagicLinkEmail(ctx, &model.SendInvitationMagicLinkEmail{
		To:        c.Email,
		URL:       url,
		FirstName: "there", // Default greeting for new users
	}); err != nil {
		return nil, err
	}

	return &dto.RequestInvitationMagicLinkOutput{
		Email: c.Email,
		IsNew: !exists,
	}, nil
}

// AuthenticateWithInvitationMagicLink authenticates a user with an invitation magic link
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
	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get organization
	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = conv.SafeValue(hostOrg.Subdomain)
	}

	// Check if user exists
	exists, err := s.Store.User().IsEmailExists(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := createRegistrationToken(c.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to generate registration token: %w", err)
		}

		return &dto.AuthenticateWithInvitationMagicLinkOutput{
			Token:     registrationToken,
			IsNewUser: true,
		}, nil
	}

	// Get existing user
	u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Create organization access
	orgAccess := &model.UserOrganizationAccess{
		ID:             uuid.Must(uuid.NewV4()),
		UserID:         u.ID,
		OrganizationID: invitedOrg.ID,
		Role:           userInvitation.Role,
	}

	// Generate token and secret
	now := time.Now()
	expiresAt := now.Add(model.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
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
		AuthURL:   config.Config.OrgBaseURL(orgSubdomain) + model.SaveAuthPath,
		Token:     token,
		Domain:    config.Config.OrgDomain(orgSubdomain),
		IsNewUser: false,
	}, nil
}

// RegisterWithInvitationMagicLink registers a new user with an invitation magic link
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
	userInvitation, err := s.Store.User().GetInvitation(ctx, storeopts.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Get organization
	invitedOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	// Verify organization access in cloud edition
	var orgSubdomain string
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		hostOrg, err := s.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
		if err != nil {
			return nil, err
		}

		if invitedOrg.ID != hostOrg.ID {
			return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
		}

		orgSubdomain = conv.SafeValue(hostOrg.Subdomain)
	}

	// Generate secret
	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	// Create new user
	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	u := &model.User{
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Email:                c.Email,
		Secret:               hashedSecret,
		EmailAuthenticatedAt: &now,
	}

	// Create organization access
	orgAccess := &model.UserOrganizationAccess{
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
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
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
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    config.Config.OrgDomain(orgSubdomain),
	}, nil
}
