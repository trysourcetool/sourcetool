package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
)

// JWTClaims is a generic constraint for all JWT claims types.
type JWTClaims interface {
	jwt.Claims
	*UserClaims | *UserEmailClaims | *UserGoogleAuthRequestClaims | *UserAuthClaims | *UserMagicLinkRegistrationClaims | *UserGoogleRegistrationClaims | *UserGoogleAuthLinkClaims
}

// NewClaims creates a new instance of the claims type.
func NewClaims[T JWTClaims]() T {
	var zero T
	switch any(zero).(type) {
	case *UserClaims:
		return any(&UserClaims{}).(T)
	case *UserEmailClaims:
		return any(&UserEmailClaims{}).(T)
	case *UserGoogleAuthRequestClaims:
		return any(&UserGoogleAuthRequestClaims{}).(T)
	case *UserAuthClaims:
		return any(&UserAuthClaims{}).(T)
	case *UserGoogleRegistrationClaims:
		return any(&UserGoogleRegistrationClaims{}).(T)
	case *UserGoogleAuthLinkClaims:
		return any(&UserGoogleAuthLinkClaims{}).(T)
	case *UserMagicLinkRegistrationClaims:
		return any(&UserMagicLinkRegistrationClaims{}).(T)
	default:
		return zero
	}
}

// SignToken is a generic function to sign JWT tokens.
func SignToken[T JWTClaims](claims T) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tok.SignedString([]byte(config.Config.Jwt.Key))
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return token, nil
}

// ParseToken is a generic function to parse JWT tokens.
func ParseToken[T JWTClaims](token string) (T, error) {
	if token == "" {
		var zero T
		return zero, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	result := NewClaims[T]()
	_, err := jwt.ParseWithClaims(token, result, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil {
		var zero T
		return zero, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return result, nil
}
