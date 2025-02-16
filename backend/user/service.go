package user

import (
	"context"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type ServiceCE interface {
	GetMe(context.Context) (*types.GetMePayload, error)
	List(context.Context) (*types.ListUsersPayload, error)
	Update(context.Context, types.UpdateUserInput) (*types.UpdateUserPayload, error)
	SendUpdateEmailInstructions(context.Context, types.SendUpdateUserEmailInstructionsInput) error
	UpdateEmail(context.Context, types.UpdateUserEmailInput) (*types.UpdateUserEmailPayload, error)
	UpdatePassword(context.Context, types.UpdateUserPasswordInput) (*types.UpdateUserPasswordPayload, error)
	SignIn(context.Context, types.SignInInput) (*types.SignInPayload, error)
	SignInWithGoogle(context.Context, types.SignInWithGoogleInput) (*types.SignInWithGooglePayload, error)
	SendSignUpInstructions(context.Context, types.SendSignUpInstructionsInput) (*types.SendSignUpInstructionsPayload, error)
	SignUp(context.Context, types.SignUpInput) (*types.SignUpPayload, error)
	SignUpWithGoogle(context.Context, types.SignUpWithGoogleInput) (*types.SignUpWithGooglePayload, error)
	RefreshToken(context.Context, types.RefreshTokenInput) (*types.RefreshTokenPayload, error)
	SaveAuth(context.Context, types.SaveAuthInput) (*types.SaveAuthPayload, error)
	ObtainAuthToken(context.Context) (*types.ObtainAuthTokenPayload, error)
	Invite(context.Context, types.InviteUsersInput) (*types.InviteUsersPayload, error)
	SignInInvitation(context.Context, types.SignInInvitationInput) (*types.SignInInvitationPayload, error)
	SignUpInvitation(context.Context, types.SignUpInvitationInput) (*types.SignUpInvitationPayload, error)
	GetGoogleAuthCodeURL(context.Context) (*types.GetGoogleAuthCodeURLPayload, error)
	GoogleOAuthCallback(context.Context, types.GoogleOAuthCallbackInput) (*types.GoogleOAuthCallbackPayload, error)
	GetGoogleAuthCodeURLInvitation(context.Context, types.GetGoogleAuthCodeURLInvitationInput) (*types.GetGoogleAuthCodeURLInvitationPayload, error)
	SignInWithGoogleInvitation(context.Context, types.SignInWithGoogleInvitationInput) (*types.SignInWithGoogleInvitationPayload, error)
	SignUpWithGoogleInvitation(context.Context, types.SignUpWithGoogleInvitationInput) (*types.SignUpWithGoogleInvitationPayload, error)
	SignOut(context.Context) (*types.SignOutPayload, error)
}

type ServiceCEImpl struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCEImpl {
	return &ServiceCEImpl{Dependency: d}
}

func (s *ServiceCEImpl) GetMe(ctx context.Context) (*types.GetMePayload, error) {
	u := ctxutils.CurrentUser(ctx)

	conds := []any{
		model.OrganizationByUserID(u.ID),
	}
	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	if subdomain != "auth" {
		conds = append(conds, model.OrganizationBySubdomain(subdomain))
	}
	o, err := s.Store.Organization().Get(ctx, conds...)
	if err != nil && !errdefs.IsOrganizationNotFound(err) {
		return nil, err
	}

	orgAccessConds := []any{
		model.UserOrganizationAccessByUserID(u.ID),
	}
	if o != nil {
		orgAccessConds = append(orgAccessConds, model.UserOrganizationAccessByOrganizationID(o.ID))
	}
	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, orgAccessConds...)
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return nil, err
	}

	userPayload := &types.UserPayload{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: strconv.FormatInt(u.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(u.UpdatedAt.Unix(), 10),
	}

	if o != nil {
		userPayload.Organization = &types.OrganizationPayload{
			ID:        o.ID.String(),
			Subdomain: o.Subdomain,
			CreatedAt: strconv.FormatInt(o.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(o.UpdatedAt.Unix(), 10),
		}
	}

	if orgAccess != nil {
		userPayload.Role = orgAccess.Role.String()
	}

	return &types.GetMePayload{
		User: userPayload,
	}, nil
}

func (s *ServiceCEImpl) List(ctx context.Context) (*types.ListUsersPayload, error) {
	o := ctxutils.CurrentOrganization(ctx)

	users, err := s.Store.User().List(ctx, model.UserByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	userInvitations, err := s.Store.User().ListInvitations(ctx, model.UserInvitationByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, model.UserOrganizationAccessByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}
	roleMap := make(map[uuid.UUID]model.UserOrganizationRole)
	for _, oa := range orgAccesses {
		roleMap[oa.UserID] = oa.Role
	}

	usersPayload := make([]*types.UserPayload, 0, len(users))
	for _, u := range users {
		usersPayload = append(usersPayload, &types.UserPayload{
			ID:        u.ID.String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Role:      roleMap[u.ID].String(),
			CreatedAt: strconv.FormatInt(u.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(u.UpdatedAt.Unix(), 10),
		})
	}

	userInvitationsPayload := make([]*types.UserInvitationPayload, 0, len(userInvitations))
	for _, ui := range userInvitations {
		userInvitationsPayload = append(userInvitationsPayload, &types.UserInvitationPayload{
			ID:        ui.ID.String(),
			Email:     ui.Email,
			CreatedAt: strconv.FormatInt(ui.CreatedAt.Unix(), 10),
		})
	}

	return &types.ListUsersPayload{
		Users:           usersPayload,
		UserInvitations: userInvitationsPayload,
	}, nil
}

func (s *ServiceCEImpl) Update(ctx context.Context, in types.UpdateUserInput) (*types.UpdateUserPayload, error) {
	currentUser := ctxutils.CurrentUser(ctx)

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

	userPayload := &types.UserPayload{
		ID:        currentUser.ID.String(),
		FirstName: currentUser.FirstName,
		LastName:  currentUser.LastName,
		Email:     currentUser.Email,
		CreatedAt: strconv.FormatInt(currentUser.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(currentUser.UpdatedAt.Unix(), 10),
	}

	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	if org != nil {
		userPayload.Organization = &types.OrganizationPayload{
			ID:        org.ID.String(),
			Subdomain: org.Subdomain,
			CreatedAt: strconv.FormatInt(org.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(org.UpdatedAt.Unix(), 10),
		}
	}

	if orgAccess != nil {
		userPayload.Role = orgAccess.Role.String()
	}

	return &types.UpdateUserPayload{
		User: userPayload,
	}, nil
}

func (s *ServiceCEImpl) SendUpdateEmailInstructions(ctx context.Context, in types.SendUpdateUserEmailInstructionsInput) error {
	if in.Email != in.EmailConfirmation {
		return errdefs.ErrInvalidArgument(errors.New("email and email confirmation do not match"))
	}

	exists, err := s.Store.User().IsEmailExists(ctx, in.Email)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(errors.New("email exists"))
	}

	currentUser := ctxutils.CurrentUser(ctx)

	tok, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID: currentUser.ID.String(),
		Email:  in.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   model.UserSignatureSubjectUpdateEmail,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.EmailTokenExpiration)),
			Issuer:    model.JwtIssuer,
		},
	})
	if err != nil {
		return err
	}

	url, err := buildUpdateEmailURL(ctx, tok)
	if err != nil {
		return err
	}

	logger.Logger.Sugar().Debug("================= URL =================")
	logger.Logger.Sugar().Debug(url)
	logger.Logger.Sugar().Debug("================= URL =================")

	if !(config.Config.Env == config.EnvLocal) {
		if err := s.Mailer.User().SendUpdateEmailInstructions(ctx, &model.SendUpdateUserEmailInstructions{
			To:        in.Email,
			FirstName: currentUser.FirstName,
			URL:       url,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceCEImpl) UpdateEmail(ctx context.Context, in types.UpdateUserEmailInput) (*types.UpdateUserEmailPayload, error) {
	c, err := s.Signer.User().ClaimsFromToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != model.UserSignatureSubjectUpdateEmail {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}
	u, err := s.Store.User().Get(ctx, model.UserByID(userID))
	if err != nil {
		return nil, err
	}

	currentUser := ctxutils.CurrentUser(ctx)
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

	userPayload := &types.UserPayload{
		ID:        currentUser.ID.String(),
		FirstName: currentUser.FirstName,
		LastName:  currentUser.LastName,
		Email:     currentUser.Email,
		CreatedAt: strconv.FormatInt(currentUser.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(currentUser.UpdatedAt.Unix(), 10),
	}

	if org != nil {
		userPayload.Organization = &types.OrganizationPayload{
			ID:        org.ID.String(),
			Subdomain: org.Subdomain,
			CreatedAt: strconv.FormatInt(org.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(org.UpdatedAt.Unix(), 10),
		}
	}

	if orgAccess != nil {
		userPayload.Role = orgAccess.Role.String()
	}

	return &types.UpdateUserEmailPayload{
		User: userPayload,
	}, nil
}

func (s *ServiceCEImpl) UpdatePassword(ctx context.Context, in types.UpdateUserPasswordInput) (*types.UpdateUserPasswordPayload, error) {
	if in.Password != in.PasswordConfirmation {
		return nil, errdefs.ErrInvalidArgument(errors.New("password and password confirmation do not match"))
	}

	currentUser := ctxutils.CurrentUser(ctx)

	h, err := hex.DecodeString(currentUser.Password)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if err = bcrypt.CompareHashAndPassword(h, []byte(in.CurrentPassword)); err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if err := model.ValidatePassword(in.Password); err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	encodedPass, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	currentUser.Password = hex.EncodeToString(encodedPass[:])

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, currentUser)
	}); err != nil {
		return nil, err
	}

	org, orgAccess, err := s.getUserOrganizationInfo(ctx)
	if err != nil {
		return nil, err
	}

	userPayload := &types.UserPayload{
		ID:        currentUser.ID.String(),
		FirstName: currentUser.FirstName,
		LastName:  currentUser.LastName,
		Email:     currentUser.Email,
		CreatedAt: strconv.FormatInt(currentUser.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(currentUser.UpdatedAt.Unix(), 10),
	}

	if org != nil {
		userPayload.Organization = &types.OrganizationPayload{
			ID:        org.ID.String(),
			Subdomain: org.Subdomain,
			CreatedAt: strconv.FormatInt(org.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(org.UpdatedAt.Unix(), 10),
		}
	}

	if orgAccess != nil {
		userPayload.Role = orgAccess.Role.String()
	}

	return &types.UpdateUserPasswordPayload{
		User: userPayload,
	}, nil
}

func (s *ServiceCEImpl) SignIn(ctx context.Context, in types.SignInInput) (*types.SignInPayload, error) {
	u, err := s.Store.User().Get(ctx, model.UserByEmail(in.Email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	h, err := hex.DecodeString(u.Password)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if err = bcrypt.CompareHashAndPassword(h, []byte(in.Password)); err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	if subdomain != "auth" {
		_, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID), model.UserOrganizationAccessByOrganizationSubdomain(subdomain))
		if err != nil {
			return nil, err
		}
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID))
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return nil, err
	}

	orgSubdomain := "auth"
	if orgAccess != nil {
		org, err := s.Store.Organization().Get(ctx, model.OrganizationByID(orgAccess.OrganizationID))
		if err != nil {
			return nil, err
		}

		orgSubdomain = org.Subdomain
	}

	now := time.Now()
	expiresAt := now.Add(model.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     in.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = secret

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &types.SignInPayload{
		AuthURL:              buildSaveAuthURL(orgSubdomain),
		Token:                token,
		Secret:               secret,
		XSRFToken:            xsrfToken,
		IsOrganizationExists: orgAccess != nil,
		Domain:               orgSubdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) SignInWithGoogle(ctx context.Context, in types.SignInWithGoogleInput) (*types.SignInWithGooglePayload, error) {
	googleAuthReqClaims, err := s.Signer.User().GoogleAuthRequestClaimsFromToken(ctx, in.SessionToken)
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

	u, err := s.Store.User().Get(ctx, model.UserByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	if subdomain != "auth" {
		_, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID), model.UserOrganizationAccessByOrganizationSubdomain(subdomain))
		if err != nil {
			return nil, err
		}
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID))
	if err != nil && !errdefs.IsUserOrganizationAccessNotFound(err) {
		return nil, err
	}

	orgSubdomain := "auth"
	if orgAccess != nil {
		org, err := s.Store.Organization().Get(ctx, model.OrganizationByID(orgAccess.OrganizationID))
		if err != nil {
			return nil, err
		}

		orgSubdomain = org.Subdomain
	}

	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     googleAuthReq.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = secret

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &types.SignInWithGooglePayload{
		AuthURL:              buildSaveAuthURL(orgSubdomain),
		Token:                token,
		Secret:               secret,
		XSRFToken:            xsrfToken,
		IsOrganizationExists: orgAccess != nil,
		Domain:               orgSubdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) SendSignUpInstructions(ctx context.Context, in types.SendSignUpInstructionsInput) (*types.SendSignUpInstructionsPayload, error) {
	exists, err := s.Store.User().IsEmailExists(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errdefs.ErrUserEmailAlreadyExists(errors.New("email exists"))
	}

	requestExists, err := s.Store.User().IsRegistrationRequestExists(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if !requestExists {
			if err := tx.User().CreateRegistrationRequest(ctx, &model.UserRegistrationRequest{
				ID:    uuid.Must(uuid.NewV4()),
				Email: in.Email,
			}); err != nil {
				return err
			}
		}

		tok, err := s.Signer.User().SignedStringFromEmail(ctx, &model.UserEmailClaims{
			Email: in.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   model.UserSignatureSubjectActivate,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.EmailTokenExpiration)),
				Issuer:    model.JwtIssuer,
			},
		})
		if err != nil {
			return err
		}

		url, err := buildUserActivateURL(ctx, tok)
		if err != nil {
			return err
		}

		logger.Logger.Sugar().Debug("================= URL =================")
		logger.Logger.Sugar().Debug(url)
		logger.Logger.Sugar().Debug("================= URL =================")

		if !(config.Config.Env == config.EnvLocal) {
			if err := s.Mailer.User().SendSignUpInstructions(ctx, &model.SendSignUpInstructions{
				To:  in.Email,
				URL: url,
			}); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &types.SendSignUpInstructionsPayload{
		Email: in.Email,
	}, nil
}

func (s *ServiceCEImpl) SignUp(ctx context.Context, in types.SignUpInput) (*types.SignUpPayload, error) {
	c, err := s.Signer.User().EmailClaimsFromToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}

	if c.Subject != model.UserSignatureSubjectActivate {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	requestUser, err := s.Store.User().GetRegistrationRequest(ctx, model.UserRegistrationRequestByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	if err := model.ValidatePassword(in.Password); err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	if in.Password != in.PasswordConfirmation {
		return nil, errdefs.ErrInvalidArgument(errors.New("password and password confirmation does not match"))
	}

	encodedPass, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	now := time.Now()
	expiresAt := now.Add(model.TmpTokenExpiration)
	u := &model.User{
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            in.FirstName,
		LastName:             in.LastName,
		Email:                c.Email,
		Password:             hex.EncodeToString(encodedPass[:]),
		Secret:               secret,
		EmailAuthenticatedAt: &now,
	}

	var token string
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		token, err = s.Signer.User().SignedString(ctx, &model.UserClaims{
			UserID:    u.ID.String(),
			Email:     c.Email,
			XSRFToken: xsrfToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				Issuer:    model.JwtIssuer,
				Subject:   model.UserSignatureSubjectEmail,
			},
		})
		if err != nil {
			return err
		}

		return tx.User().DeleteRegistrationRequest(ctx, requestUser)
	}); err != nil {
		return nil, err
	}

	if !(config.Config.Env == config.EnvLocal) {
		s.Mailer.User().SendWelcomeEmail(ctx, &model.SendWelcomeEmail{
			To: u.Email,
		})
	}

	return &types.SignUpPayload{
		Token:     token,
		XSRFToken: xsrfToken,
	}, nil
}

func (s *ServiceCEImpl) SignUpWithGoogle(ctx context.Context, in types.SignUpWithGoogleInput) (*types.SignUpWithGooglePayload, error) {
	googleAuthReqClaims, err := s.Signer.User().GoogleAuthRequestClaimsFromToken(ctx, in.SessionToken)
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

	requestUser, err := s.Store.User().GetRegistrationRequest(ctx, model.UserRegistrationRequestByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	secret, err := generateSecret()
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
		Secret:               secret,
		EmailAuthenticatedAt: &now,
		GoogleID:             googleAuthReq.GoogleID,
	}

	var token string
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		token, err = s.Signer.User().SignedString(ctx, &model.UserClaims{
			UserID:    u.ID.String(),
			Email:     googleAuthReq.Email,
			XSRFToken: xsrfToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				Issuer:    model.JwtIssuer,
				Subject:   model.UserSignatureSubjectEmail,
			},
		})
		if err != nil {
			return err
		}

		return tx.User().DeleteRegistrationRequest(ctx, requestUser)
	}); err != nil {
		return nil, err
	}

	return &types.SignUpWithGooglePayload{
		Token:     token,
		XSRFToken: xsrfToken,
	}, nil
}

func (s *ServiceCEImpl) RefreshToken(ctx context.Context, in types.RefreshTokenInput) (*types.RefreshTokenPayload, error) {
	if in.XSRFTokenCookie != in.XSRFTokenHeader {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid xsrf token"))
	}

	u, err := s.Store.User().Get(ctx, model.UserBySecret(in.Secret))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	_, err = s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID), model.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     u.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	return &types.RefreshTokenPayload{
		Token:     token,
		Secret:    u.Secret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) SaveAuth(ctx context.Context, in types.SaveAuthInput) (*types.SaveAuthPayload, error) {
	c, err := s.Signer.User().ClaimsFromToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}

	u, err := s.Store.User().Get(ctx, model.UserByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	_, err = s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID), model.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(model.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     u.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = secret

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &types.SaveAuthPayload{
		Token:       token,
		Secret:      secret,
		XSRFToken:   xsrfToken,
		ExpiresAt:   strconv.FormatInt(expiresAt.Unix(), 10),
		RedirectURL: buildServiceURL(subdomain),
		Domain:      subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) ObtainAuthToken(ctx context.Context) (*types.ObtainAuthTokenPayload, error) {
	u := ctxutils.CurrentUser(ctx)

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, err
	}

	o, err := s.Store.Organization().Get(ctx, model.OrganizationByID(orgAccess.OrganizationID))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(model.TmpTokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     u.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	if err = s.Store.RunTransaction(func(tx infra.Transaction) error {
		return tx.User().Update(ctx, u)
	}); err != nil {
		return nil, err
	}

	return &types.ObtainAuthTokenPayload{
		AuthURL: buildSaveAuthURL(o.Subdomain),
		Token:   token,
	}, nil
}

func (s *ServiceCEImpl) Invite(ctx context.Context, in types.InviteUsersInput) (*types.InviteUsersPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditUser); err != nil {
		return nil, err
	}

	o := ctxutils.CurrentOrganization(ctx)
	u := ctxutils.CurrentUser(ctx)

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

		tok, err := s.Signer.User().SignedStringFromEmail(ctx, &model.UserEmailClaims{
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   model.UserSignatureSubjectInvitation,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.EmailTokenExpiration)),
				Issuer:    model.JwtIssuer,
			},
		})
		if err != nil {
			return nil, err
		}

		url, err := buildInvitationURL(ctx, tok, email, userExists)
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

	usersInvitationsPayload := make([]*types.UserInvitationPayload, 0, len(invitations))
	for _, ui := range invitations {
		usersInvitationsPayload = append(usersInvitationsPayload, &types.UserInvitationPayload{
			ID:        ui.ID.String(),
			Email:     ui.Email,
			CreatedAt: strconv.FormatInt(ui.CreatedAt.Unix(), 10),
		})
	}

	return &types.InviteUsersPayload{
		UserInvitations: usersInvitationsPayload,
	}, nil
}

func (s *ServiceCEImpl) SignInInvitation(ctx context.Context, in types.SignInInvitationInput) (*types.SignInInvitationPayload, error) {
	c, err := s.Signer.User().EmailClaimsFromToken(ctx, in.InvitationToken)
	if err != nil {
		return nil, err
	}

	if c.Subject != model.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, model.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, model.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	hostOrg, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	if invitedOrg.ID != hostOrg.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	u, err := s.Store.User().Get(ctx, model.UserByEmail(c.Email))
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
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     u.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = secret

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

	return &types.SignInInvitationPayload{
		Token:     token,
		Secret:    secret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) SignUpInvitation(ctx context.Context, in types.SignUpInvitationInput) (*types.SignUpInvitationPayload, error) {
	c, err := s.Signer.User().EmailClaimsFromToken(ctx, in.InvitationToken)
	if err != nil {
		return nil, err
	}

	if c.Subject != model.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, model.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, model.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	hostOrg, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	if invitedOrg.ID != hostOrg.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	if err := model.ValidatePassword(in.Password); err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	if in.Password != in.PasswordConfirmation {
		return nil, errdefs.ErrInvalidArgument(errors.New("password and password confirmation does not match"))
	}

	encodedPass, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	secret, err := generateSecret()
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
		Password:             hex.EncodeToString(encodedPass[:]),
		Secret:               secret,
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

		token, err = s.Signer.User().SignedString(ctx, &model.UserClaims{
			UserID:    u.ID.String(),
			Email:     c.Email,
			XSRFToken: xsrfToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				Issuer:    model.JwtIssuer,
				Subject:   model.UserSignatureSubjectEmail,
			},
		})
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &types.SignUpInvitationPayload{
		Token:     token,
		Secret:    secret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) GetGoogleAuthCodeURL(ctx context.Context) (*types.GetGoogleAuthCodeURLPayload, error) {
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
			Domain:    buildServiceDomain("auth"),
			ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour),
			Invited:   false,
		})
	}); err != nil {
		return nil, err
	}

	return &types.GetGoogleAuthCodeURLPayload{
		URL: url,
	}, nil
}

func (s *ServiceCEImpl) GoogleOAuthCallback(ctx context.Context, in types.GoogleOAuthCallbackInput) (*types.GoogleOAuthCallbackPayload, error) {
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

	sessionToken, err := s.Signer.User().SignedStringGoogleAuthRequest(ctx, &model.UserGoogleAuthRequestClaims{
		GoogleAuthRequestID: googleAuthReq.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(googleAuthReq.ExpiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectGoogleAuthRequest,
		},
	})
	if err != nil {
		return nil, err
	}

	if isUserExists {
		u, err := s.Store.User().Get(ctx, model.UserByEmail(userInfo.email))
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

		return &types.GoogleOAuthCallbackPayload{
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

	return &types.GoogleOAuthCallbackPayload{
		SessionToken: sessionToken,
		IsUserExists: isUserExists,
		FirstName:    userInfo.givenName,
		LastName:     userInfo.familyName,
		Domain:       googleAuthReq.Domain,
		Invited:      googleAuthReq.Invited,
	}, nil
}

func (s *ServiceCEImpl) GetGoogleAuthCodeURLInvitation(ctx context.Context, in types.GetGoogleAuthCodeURLInvitationInput) (*types.GetGoogleAuthCodeURLInvitationPayload, error) {
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

	c, err := s.Signer.User().EmailClaimsFromToken(ctx, in.InvitationToken)
	if err != nil {
		return nil, err
	}
	if c.Subject != model.UserSignatureSubjectInvitation {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid jwt subject"))
	}

	userInvitation, err := s.Store.User().GetInvitation(ctx, model.UserInvitationByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, model.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	hostOrg, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	if invitedOrg.ID != hostOrg.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.User().BulkDeleteGoogleAuthRequests(ctx, googleAuthReqs); err != nil {
			return err
		}

		return tx.User().CreateGoogleAuthRequest(ctx, &model.UserGoogleAuthRequest{
			ID:        state,
			Domain:    ctxutils.HTTPHost(ctx),
			ExpiresAt: time.Now().Add(time.Duration(24) * time.Hour),
			Invited:   true,
		})
	}); err != nil {
		return nil, err
	}

	return &types.GetGoogleAuthCodeURLInvitationPayload{
		URL: url,
	}, nil
}

func (s *ServiceCEImpl) SignInWithGoogleInvitation(ctx context.Context, in types.SignInWithGoogleInvitationInput) (*types.SignInWithGoogleInvitationPayload, error) {
	googleAuthReqClaims, err := s.Signer.User().GoogleAuthRequestClaimsFromToken(ctx, in.SessionToken)
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

	userInvitation, err := s.Store.User().GetInvitation(ctx, model.UserInvitationByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, model.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	hostOrg, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	if invitedOrg.ID != hostOrg.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	u, err := s.Store.User().Get(ctx, model.UserByEmail(googleAuthReq.Email))
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
	token, err := s.Signer.User().SignedString(ctx, &model.UserClaims{
		UserID:    u.ID.String(),
		Email:     googleAuthReq.Email,
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    model.JwtIssuer,
			Subject:   model.UserSignatureSubjectEmail,
		},
	})
	if err != nil {
		return nil, err
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, errdefs.ErrInternal(err)
	}

	u.Secret = secret

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

	return &types.SignInWithGoogleInvitationPayload{
		Token:     token,
		Secret:    secret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) SignUpWithGoogleInvitation(ctx context.Context, in types.SignUpWithGoogleInvitationInput) (*types.SignUpWithGoogleInvitationPayload, error) {
	googleAuthReqClaims, err := s.Signer.User().GoogleAuthRequestClaimsFromToken(ctx, in.SessionToken)
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

	userInvitation, err := s.Store.User().GetInvitation(ctx, model.UserInvitationByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	invitedOrg, err := s.Store.Organization().Get(ctx, model.OrganizationByID(userInvitation.OrganizationID))
	if err != nil {
		return nil, err
	}

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	hostOrg, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	if invitedOrg.ID != hostOrg.ID {
		return nil, errdefs.ErrUnauthenticated(errors.New("invalid organization"))
	}

	if err := s.Store.User().DeleteGoogleAuthRequest(ctx, googleAuthReq); err != nil {
		return nil, err
	}

	if time.Now().After(googleAuthReq.ExpiresAt) {
		return nil, errdefs.ErrInvalidArgument(errors.New("google auth code expired"))
	}

	requestUser, err := s.Store.User().GetRegistrationRequest(ctx, model.UserRegistrationRequestByEmail(googleAuthReq.Email))
	if err != nil {
		return nil, err
	}

	secret, err := generateSecret()
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
		Secret:               secret,
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

		token, err = s.Signer.User().SignedString(ctx, &model.UserClaims{
			UserID:    u.ID.String(),
			Email:     googleAuthReq.Email,
			XSRFToken: xsrfToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				Issuer:    model.JwtIssuer,
				Subject:   model.UserSignatureSubjectEmail,
			},
		})
		if err != nil {
			return err
		}

		return tx.User().DeleteRegistrationRequest(ctx, requestUser)
	}); err != nil {
		return nil, err
	}

	return &types.SignUpWithGoogleInvitationPayload{
		Token:     token,
		Secret:    secret,
		XSRFToken: xsrfToken,
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		Domain:    subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) SignOut(ctx context.Context) (*types.SignOutPayload, error) {
	u := ctxutils.CurrentUser(ctx)

	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]
	_, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID), model.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	return &types.SignOutPayload{
		Domain: subdomain + "." + config.Config.Domain,
	}, nil
}

func (s *ServiceCEImpl) createPersonalAPIKey(ctx context.Context, tx infra.Transaction, u *model.User, org *model.Organization) error {
	devEnv, err := s.Store.Environment().Get(ctx, model.EnvironmentByOrganizationID(org.ID), model.EnvironmentBySlug(model.EnvironmentSlugDevelopment))
	if err != nil {
		return err
	}

	key, err := model.GenerateAPIKey(org.Subdomain, devEnv.Slug)
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

func (s *ServiceCEImpl) getUserOrganizationInfo(ctx context.Context) (*model.Organization, *model.UserOrganizationAccess, error) {
	u := ctxutils.CurrentUser(ctx)
	subdomain := strings.Split(ctxutils.HTTPHost(ctx), ".")[0]

	o, err := s.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, nil, err
	}

	orgAccess, err := s.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByOrganizationID(o.ID), model.UserOrganizationAccessByUserID(u.ID))
	if err != nil {
		return nil, nil, err
	}

	return o, orgAccess, nil
}
