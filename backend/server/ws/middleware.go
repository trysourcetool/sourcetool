package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/jwt"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type Middleware interface {
	Auth(next http.Handler) http.Handler
}

type MiddlewareCE struct {
	infra.Store
}

func NewMiddlewareCE(s infra.Store) *MiddlewareCE {
	return &MiddlewareCE{s}
}

type ClientHeader struct {
	Token string `json:"Authorization"`
}

func (m *MiddlewareCE) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		subdomain, err := m.getSubdomainIfCloudEdition(r)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)

		if token, err := r.Cookie("access_token"); err == nil {
			u, err := m.getCurrentUser(ctx, r, token.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)
		} else if apiKeyHeader := r.Header.Get("Authorization"); apiKeyHeader != "" {
			apikeyVal, err := m.extractIncomingToken(apiKeyHeader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			apikey, err := m.Store.APIKey().Get(ctx, storeopts.APIKeyByKey(apikeyVal))
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

func (m *MiddlewareCE) getCurrentUser(ctx context.Context, r *http.Request, token string) (*model.User, error) {
	c, err := m.validateUserToken(token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	u, err := m.Store.User().Get(ctx, storeopts.UserByID(userID))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *MiddlewareCE) getCurrentOrganization(ctx context.Context, subdomain string) (*model.Organization, error) {
	opts := []storeopts.OrganizationOption{}
	if subdomain != "" && subdomain != "auth" {
		opts = append(opts, storeopts.OrganizationBySubdomain(subdomain))
	}

	return m.Store.Organization().Get(ctx, opts...)
}

func (m *MiddlewareCE) validateUserToken(token string) (*jwt.UserAuthClaims, error) {
	if token == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get token"))
	}

	claims, err := jwt.ParseToken[*jwt.UserAuthClaims](token)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return claims, nil
}

func (m *MiddlewareCE) extractIncomingToken(headerValue string) (string, error) {
	if !strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
		return "", fmt.Errorf("invalid or malformed %q header, expected 'Bearer JWT-token...'", headerValue)
	}
	return strings.Split(headerValue, " ")[1], nil
}

func (m *MiddlewareCE) getSubdomainIfCloudEdition(r *http.Request) (string, error) {
	if !config.Config.IsCloudEdition {
		return "", nil
	}
	return httputil.GetSubdomainFromHost(r.Host)
}
