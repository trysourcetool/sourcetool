package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func signToken(claims jwt.Claims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tok.SignedString([]byte(config.Config.Jwt.Key))
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return token, nil
}

func SignAuthToken(userID, xsrfToken string, expiresAt time.Time) (string, error) {
	return signToken(&AuthClaims{
		XSRFToken: xsrfToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    issuer,
			Subject:   userID,
		},
	})
}

func ParseAuthClaims(token string) (*AuthClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &AuthClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignUpdateUserEmailToken(userID, email string) (string, error) {
	return signToken(&UpdateUserEmailClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    issuer,
			Subject:   userID,
		},
	})
}

func ParseUpdateUserEmailClaims(token string) (*UpdateUserEmailClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &UpdateUserEmailClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignInvitationToken(email string) (string, error) {
	return signToken(&InvitationClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    issuer,
			Subject:   email,
		},
	})
}

func ParseInvitationClaims(token string) (*InvitationClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &InvitationClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignMagicLinkToken(email string) (string, error) {
	return signToken(&MagicLinkClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    issuer,
			Subject:   email,
		},
	})
}

func ParseMagicLinkClaims(token string) (*MagicLinkClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &MagicLinkClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignInvitationMagicLinkToken(email string) (string, error) {
	return signToken(&InvitationMagicLinkClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    issuer,
			Subject:   email,
		},
	})
}

func ParseInvitationMagicLinkClaims(token string) (*InvitationMagicLinkClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &InvitationMagicLinkClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignMagicLinkRegistrationToken(email string) (string, error) {
	return signToken(&MagicLinkRegistrationClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    issuer,
			Subject:   email,
		},
	})
}

func ParseMagicLinkRegistrationClaims(token string) (*MagicLinkRegistrationClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &MagicLinkRegistrationClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignGoogleAuthLinkToken(flow GoogleAuthFlow, invitationOrgID uuid.UUID, hostSubdomain string) (string, error) {
	return signToken(&GoogleAuthLinkClaims{
		Flow:            flow,
		InvitationOrgID: invitationOrgID,
		HostSubdomain:   hostSubdomain,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    issuer,
		},
	})
}

func ParseGoogleAuthLinkClaims(token string) (*GoogleAuthLinkClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &GoogleAuthLinkClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func SignGoogleRegistrationToken(googleID, email, firstName, lastName string, flow GoogleAuthFlow, invitationOrgID uuid.UUID, role string) (string, error) {
	return signToken(&GoogleRegistrationClaims{
		GoogleID:        googleID,
		FirstName:       firstName,
		LastName:        lastName,
		Flow:            flow,
		InvitationOrgID: invitationOrgID,
		Role:            role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    issuer,
			Subject:   email,
		},
	})
}

func ParseGoogleRegistrationClaims(token string) (*GoogleRegistrationClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &GoogleRegistrationClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}
