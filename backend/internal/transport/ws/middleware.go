package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/apikey"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/httpx"
)

type Middleware interface {
	Auth(next http.Handler) http.Handler
}

type MiddlewareCE struct {
	db.Repository
}

func NewMiddlewareCE(r db.Repository) *MiddlewareCE {
	return &MiddlewareCE{r}
}

type ClientHeader struct {
	Token string `json:"Authorization"`
}

func (m *MiddlewareCE) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		subdomain, err := m.getSubdomainIfCloudEdition(r)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
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

			apikey, err := m.Repository.APIKey().Get(ctx, apikey.ByKey(apikeyVal))
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

func (m *MiddlewareCE) getCurrentUser(ctx context.Context, r *http.Request, token string) (*user.User, error) {
	c, err := m.validateUserToken(token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	u, err := m.Repository.User().Get(ctx, user.ByID(userID))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *MiddlewareCE) getCurrentOrganization(ctx context.Context, subdomain string) (*organization.Organization, error) {
	opts := []organization.RepositoryOption{}
	if subdomain != "" && subdomain != "auth" {
		opts = append(opts, organization.BySubdomain(subdomain))
	}

	return m.Repository.Organization().Get(ctx, opts...)
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
	return httpx.GetSubdomainFromHost(r.Host)
}
