package user

import (
	"context"
	"errors"
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
	"github.com/trysourcetool/sourcetool/backend/logger"
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
	UpdatePassword(context.Context, dto.UpdateUserPasswordInput) (*dto.UpdateUserPasswordOutput, error)

	// Authentication methods
	SignIn(context.Context, dto.SignInInput) (*dto.SignInOutput, error)
	SignOut(context.Context) (*dto.SignOutOutput, error)
	RefreshToken(context.Context, dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error)
	SaveAuth(context.Context, dto.SaveAuthInput) (*dto.SaveAuthOutput, error)
	ObtainAuthToken(context.Context) (*dto.ObtainAuthTokenOutput, error)

	// Registration methods
	SendSignUpInstructions(context.Context, dto.SendSignUpInstructionsInput) (*dto.SendSignUpInstructionsOutput, error)
	SignUp(context.Context, dto.SignUpInput) (*dto.SignUpOutput, error)

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

	// Send email with URL (or log it in local environment)
	return s.sendEmailWithLogging(ctx, url, func() error {
		return s.Mailer.User().SendUpdateEmailInstructions(ctx, &model.SendUpdateUserEmailInstructions{
			To:        in.Email,
			FirstName: currentUser.FirstName,
			URL:       url,
		})
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

// UpdatePassword changes the user's password after verifying their current password.
func (s *ServiceCE) UpdatePassword(ctx context.Context, in dto.UpdateUserPasswordInput) (*dto.UpdateUserPasswordOutput, error) {
	currentUser := ctxutil.CurrentUser(ctx)

	// Verify current password
	if err := verifyPassword(currentUser.Password, in.CurrentPassword); err != nil {
		return nil, err
	}

	// Hash new password
	hashedPassword, err := hashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	// Update user record
	currentUser.Password = hashedPassword
	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, currentUser)
	}); err != nil {
		return nil, err
	}

	// Get organization information for response
	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	var role model.UserOrganizationRole
	if orgAccess != nil {
		role = orgAccess.Role
	}

	return &dto.UpdateUserPasswordOutput{
		User: dto.UserFromModel(currentUser, org, role),
	}, nil
}

// SignIn authenticates a user with email and password.
func (s *ServiceCE) SignIn(ctx context.Context, in dto.SignInInput) (*dto.SignInOutput, error) {
	// Get user by email
	u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(in.Email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	// Verify password
	if err = verifyPassword(u.Password, in.Password); err != nil {
		return nil, err
	}

	// Get user's organization access information
	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}

	// Handle organization subdomain logic
	subdomain := ctxutil.Subdomain(ctx)
	var orgAccess *model.UserOrganizationAccess
	var orgSubdomain string

	if config.Config.IsCloudEdition {
		if subdomain != "auth" {
			// For specific organization subdomain
			orgAccess, err = s.Store.User().GetOrganizationAccess(ctx,
				storeopts.UserOrganizationAccessByUserID(u.ID),
				storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
			if err != nil {
				return nil, err
			}
			orgSubdomain = subdomain
		} else {
			// For auth subdomain
			if len(orgAccesses) == 1 {
				// Single organization - redirect to it
				orgAccess = orgAccesses[0]
				o, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
				if err != nil {
					return nil, err
				}
				orgSubdomain = conv.SafeValue(o.Subdomain)
			} else {
				// Multiple organizations - send email with options
				return s.handleMultipleOrganizations(ctx, u, orgAccesses)
			}
		}
	} else {
		// Self-hosted mode has only one organization
		orgAccess = orgAccesses[0]
	}

	// Create token, secret, etc.
	token, xsrfToken, plainSecret, hashedSecret, _, err := s.createTokenWithSecret(
		u.ID, model.TmpTokenExpiration)
	if err != nil {
		return nil, err
	}

	// Update user with new secret
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

	return &dto.SignInOutput{
		AuthURL:              authURL,
		Token:                token,
		Secret:               plainSecret,
		XSRFToken:            xsrfToken,
		IsOrganizationExists: orgAccess != nil,
		Domain:               config.Config.OrgDomain(orgSubdomain),
	}, nil
}

// handleMultipleOrganizations sends an email with login links when a user belongs to multiple organizations.
func (s *ServiceCE) handleMultipleOrganizations(ctx context.Context, u *model.User, orgAccesses []*model.UserOrganizationAccess) (*dto.SignInOutput, error) {
	loginURLs := make([]string, 0, len(orgAccesses))

	for _, access := range orgAccesses {
		org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(access.OrganizationID))
		if err != nil {
			return nil, err
		}
		loginURLs = append(loginURLs, config.Config.OrgBaseURL(conv.SafeValue(org.Subdomain))+"/signin")
	}

	if err := s.Mailer.User().SendMultipleOrganizationsEmail(ctx, &model.SendMultipleOrganizationsEmail{
		To:        u.Email,
		FirstName: u.FirstName,
		Email:     u.Email,
		LoginURLs: loginURLs,
	}); err != nil {
		return nil, err
	}

	return nil, errdefs.ErrUserMultipleOrganizations(errors.New("email belongs to multiple organizations"))
}

func (s *ServiceCE) SignInWithGoogle(ctx context.Context, in dto.SignInWithGoogleInput) (*dto.SignInWithGoogleOutput, error) {
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

	u, err := s.Store.User().Get(ctx, storeopts.UserByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		if subdomain != "auth" {
			_, err := s.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(u.ID), storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
			if err != nil {
				return nil, err
			}
		}
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(u.ID))
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return nil, err
	}

	orgSubdomain := "auth"
	if orgAccess != nil {
		org, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
		if err != nil {
			return nil, err
		}

		orgSubdomain = conv.SafeValue(org.Subdomain)
	}

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

	u.Secret = hashedSecret

	authURL, err := buildSaveAuthURL(orgSubdomain)
	if err != nil {
		return nil, err
	}

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

// SendSignUpInstructions sends an email with a sign-up activation link.
func (s *ServiceCE) SendSignUpInstructions(ctx context.Context, in dto.SendSignUpInstructionsInput) (*dto.SendSignUpInstructionsOutput, error) {
	// Check self-hosted organization restriction
	if err := s.validateSelfHostedOrganization(ctx); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.Store.User().IsEmailExists(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errdefs.ErrUserEmailAlreadyExists(errors.New("email already exists"))
	}

	// In staging environment, only allow @trysourcetool.com email addresses
	if config.Config.Env == config.EnvStaging && !strings.HasSuffix(in.Email, "@trysourcetool.com") {
		return &dto.SendSignUpInstructionsOutput{
			Email: in.Email,
		}, nil
	}

	// Check if registration request already exists
	requestExists, err := s.Store.User().IsRegistrationRequestExists(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	// Create or update registration request and send instructions
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		// Create registration request if needed
		if !requestExists {
			if err := tx.User().CreateRegistrationRequest(ctx, &model.UserRegistrationRequest{
				ID:    uuid.Must(uuid.NewV4()),
				Email: in.Email,
			}); err != nil {
				return err
			}
		}

		// Create token for activation
		tok, err := createUserToken("", in.Email, time.Now().Add(model.EmailTokenExpiration), jwt.UserSignatureSubjectActivate)
		if err != nil {
			return err
		}

		// Build activation URL
		url, err := buildUserActivateURL(tok)
		if err != nil {
			return err
		}

		// Send email with URL (or log it in local environment)
		return s.sendEmailWithLogging(ctx, url, func() error {
			return s.Mailer.User().SendSignUpInstructions(ctx, &model.SendSignUpInstructions{
				To:  in.Email,
				URL: url,
			})
		})
	}); err != nil {
		return nil, err
	}

	return &dto.SendSignUpInstructionsOutput{
		Email: in.Email,
	}, nil
}

// SignUp registers a new user with the activated token.
func (s *ServiceCE) SignUp(ctx context.Context, in dto.SignUpInput) (*dto.SignUpOutput, error) {
	// Verify token
	c, err := jwt.ParseToken[*jwt.UserEmailClaims](in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != jwt.UserSignatureSubjectActivate {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	// Get registration request
	requestUser, err := s.Store.User().GetRegistrationRequest(ctx, storeopts.UserRegistrationRequestByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	// Create password hash
	hashedPassword, err := hashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	// Generate user secret
	plainSecret, hashedSecret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	// Prepare user object
	now := time.Now()
	u := &model.User{
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Email:                c.Email,
		Password:             hashedPassword,
		Secret:               hashedSecret,
		EmailAuthenticatedAt: &now,
	}

	// Prepare token expiration time based on edition
	var token string
	var xsrfToken string

	// Create user and setup related data in transaction
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		// Handle different editions
		expiration := model.TmpTokenExpiration
		if !config.Config.IsCloudEdition {
			// For self-hosted, create initial organization
			if err := s.createInitialOrganizationForSelfHosted(ctx, tx, u); err != nil {
				return err
			}
			expiration = model.TokenExpiration()
		}

		// Create token
		var err error
		token, xsrfToken, _, _, _, err = s.createTokenWithSecret(u.ID, expiration)
		if err != nil {
			return err
		}

		// Delete registration request
		return tx.User().DeleteRegistrationRequest(ctx, requestUser)
	}); err != nil {
		return nil, err
	}

	return &dto.SignUpOutput{
		Token:     token,
		Secret:    plainSecret,
		XSRFToken: xsrfToken,
	}, nil
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
	if in.XSRFTokenCookie != in.XSRFTokenHeader {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid xsrf token"))
	}

	hashedSecret := hashSecret(in.Secret)
	u, err := s.Store.User().Get(ctx, storeopts.UserBySecret(hashedSecret))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		_, err = s.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(u.ID), storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
		if err != nil {
			return nil, err
		}
	}

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
		Domain:    config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) SaveAuth(ctx context.Context, in dto.SaveAuthInput) (*dto.SaveAuthOutput, error) {
	c, err := jwt.ParseToken[*jwt.UserAuthClaims](in.Token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, err
	}

	u, err := s.Store.User().Get(ctx, storeopts.UserByID(userID))
	if err != nil {
		return nil, err
	}

	var subdomain string
	if config.Config.IsCloudEdition {
		subdomain = ctxutil.Subdomain(ctx)
		_, err = s.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(u.ID), storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
		if err != nil {
			return nil, err
		}
	}

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

	u.Secret = hashedSecret

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
		RedirectURL: config.Config.OrgBaseURL(subdomain),
		Domain:      config.Config.OrgDomain(subdomain),
	}, nil
}

func (s *ServiceCE) ObtainAuthToken(ctx context.Context) (*dto.ObtainAuthTokenOutput, error) {
	u := ctxutil.CurrentUser(ctx)

	orgAccess, err := s.Store.User().GetOrganizationAccess(
		ctx,
		storeopts.UserOrganizationAccessByUserID(u.ID),
		storeopts.UserOrganizationAccessOrderBy("created_at DESC"),
	)
	if err != nil {
		return nil, err
	}

	o, err := s.Store.Organization().Get(ctx, storeopts.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(model.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := createAuthToken(u.ID.String(), xsrfToken, expiresAt, jwt.UserSignatureSubjectEmail)
	if err != nil {
		return nil, err
	}

	authURL, err := buildSaveAuthURL(conv.SafeValue(o.Subdomain))
	if err != nil {
		return nil, err
	}

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

		userExists, err := s.Store.User().IsEmailExists(ctx, email)
		if err != nil {
			return nil, err
		}

		tok, err := createUserToken("", email, time.Now().Add(model.EmailTokenExpiration), jwt.UserSignatureSubjectInvitation)
		if err != nil {
			return nil, err
		}

		url, err := buildInvitationURL(conv.SafeValue(o.Subdomain), tok, email, userExists)
		if err != nil {
			return nil, err
		}

		emailInput.URLs[email] = url

		logger.Logger.Sugar().Debug("================= URL =================")
		logger.Logger.Sugar().Debug(url)
		logger.Logger.Sugar().Debug("================= URL =================")

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
		GoogleID:             googleAuthReq.GoogleID,
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

	userExists, err := s.Store.User().IsEmailExists(ctx, userInvitation.Email)
	if err != nil {
		return nil, err
	}

	tok, err := createUserEmailToken(userInvitation.Email, time.Now().Add(model.EmailTokenExpiration), jwt.UserSignatureSubjectInvitation)
	if err != nil {
		return nil, err
	}

	url, err := buildInvitationURL(conv.SafeValue(o.Subdomain), tok, userInvitation.Email, userExists)
	if err != nil {
		return nil, err
	}

	emailInput := &model.SendInvitationEmail{
		Invitees: u.FullName(),
		URLs:     map[string]string{userInvitation.Email: url},
	}

	logger.Logger.Sugar().Debug("================= URL =================")
	logger.Logger.Sugar().Debug(url)
	logger.Logger.Sugar().Debug("================= URL =================")

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

// sendEmailWithLogging sends an email with logging for local development.
func (s *ServiceCE) sendEmailWithLogging(ctx context.Context, url string, sendFunc func() error) error {
	logger.Logger.Sugar().Debug("================= URL =================")
	logger.Logger.Sugar().Debug(url)
	logger.Logger.Sugar().Debug("================= URL =================")

	if !(config.Config.Env == config.EnvLocal) {
		if err := sendFunc(); err != nil {
			return err
		}
	}

	return nil
}

// getUserOrganizationInfo retrieves organization and access information for the current user.
func (s *ServiceCE) getUserOrganizationInfo(ctx context.Context) (*model.Organization, *model.UserOrganizationAccess, error) {
	u := ctxutil.CurrentUser(ctx)

	// Build organization options based on edition
	orgOpts := []storeopts.OrganizationOption{
		storeopts.OrganizationByUserID(u.ID),
	}

	// Add subdomain filter in cloud edition
	if config.Config.IsCloudEdition {
		subdomain := ctxutil.Subdomain(ctx)
		orgOpts = append(orgOpts, storeopts.OrganizationBySubdomain(subdomain))
	}

	// Get organization
	o, err := s.Store.Organization().Get(ctx, orgOpts...)
	if err != nil {
		return nil, nil, err
	}

	// Get user's access role in this organization
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx,
		storeopts.UserOrganizationAccessByOrganizationID(o.ID),
		storeopts.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return o, orgAccess, nil
}
