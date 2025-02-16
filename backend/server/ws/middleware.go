package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type MiddlewareCE interface {
	Auth(next http.Handler) http.Handler
}

type MiddlewareCEImpl struct {
	infra.Store
}

func NewMiddlewareCE(s infra.Store) *MiddlewareCEImpl {
	return &MiddlewareCEImpl{s}
}

type ClientHeader struct {
	Token string `json:"Authorization"`
}

func (m *MiddlewareCEImpl) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		subdomain := strings.Split(r.Host, ".")[0]
		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutils.CurrentOrganizationCtxKey, o)

		if token, err := r.Cookie("access_token"); err == nil {
			u, err := m.getCurrentUser(ctx, r, token.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, ctxutils.CurrentUserCtxKey, u)
		} else if apiKeyHeader := r.Header.Get("Authorization"); apiKeyHeader != "" {
			apikeyVal, err := m.extractIncomingToken(apiKeyHeader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			apikey, err := m.Store.APIKey().Get(ctx, model.APIKeyByKey(apikeyVal))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if o.ID != apikey.OrganizationID {
				http.Error(w, "organization not found", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "failed to get token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCEImpl) getCurrentUser(ctx context.Context, r *http.Request, token string) (*model.User, error) {
	c, err := m.validateUserToken(token)
	if err != nil {
		return nil, err
	}

	u, err := m.Store.User().Get(ctx, model.UserByEmail(c.Email))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *MiddlewareCEImpl) getCurrentOrganization(ctx context.Context, subdomain string) (*model.Organization, error) {
	o, err := m.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (m *MiddlewareCEImpl) validateUserToken(token string) (*model.UserClaims, error) {
	if token == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get token"))
	}

	claims := &model.UserClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.Jwt.Key), nil
	})
	if err != nil || !tok.Valid {
		return nil, errdefs.ErrUnauthenticated(fmt.Errorf("failed to parse token: %s", err))
	}

	return claims, nil
}

func (m *MiddlewareCEImpl) extractIncomingToken(headerValue string) (string, error) {
	if !strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
		return "", fmt.Errorf("invalid or malformed %q header, expected 'Bearer JWT-token...'", headerValue)
	}
	return strings.Split(headerValue, " ")[1], nil
}
