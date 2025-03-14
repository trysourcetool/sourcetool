package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type SignerCE struct{}

func NewSignerCE() *SignerCE {
	return &SignerCE{}
}

func (s *SignerCE) SignedString(ctx context.Context, e *model.UserClaims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, e)

	token, err := tok.SignedString([]byte(config.Config.Jwt.Key))
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return token, nil
}

func (s *SignerCE) SignedStringFromEmail(ctx context.Context, e *model.UserEmailClaims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, e)

	token, err := tok.SignedString([]byte(config.Config.Jwt.Key))
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return token, nil
}

func (s *SignerCE) SignedStringGoogleAuthRequest(ctx context.Context, e *model.UserGoogleAuthRequestClaims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, e)

	token, err := tok.SignedString([]byte(config.Config.Jwt.Key))
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return token, nil
}

func (s *SignerCE) ClaimsFromToken(ctx context.Context, token string) (*model.UserClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &model.UserClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil || !tok.Valid {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func (s *SignerCE) EmailClaimsFromToken(ctx context.Context, token string) (*model.UserEmailClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &model.UserEmailClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil || !tok.Valid {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func (s *SignerCE) GoogleAuthRequestClaimsFromToken(ctx context.Context, token string) (*model.UserGoogleAuthRequestClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &model.UserGoogleAuthRequestClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil || !tok.Valid {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func (s *SignerCE) AuthClaimsFromToken(ctx context.Context, token string) (*model.UserAuthClaims, error) {
	if token == "" {
		return nil, errdefs.ErrInternal(errors.New("failed to get token"))
	}

	claims := &model.UserAuthClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil || !tok.Valid {
		return nil, errdefs.ErrInternal(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func (s *SignerCE) SignedStringAuth(ctx context.Context, e *model.UserAuthClaims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, e)

	token, err := tok.SignedString([]byte(config.Config.Jwt.Key))
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return token, nil
}
