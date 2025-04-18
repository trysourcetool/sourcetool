package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/organization"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/httpx"
)

type Middleware interface {
	AuthUser(next http.Handler) http.Handler
	AuthUserWithOrganization(next http.Handler) http.Handler
	AuthUserWithOrganizationIfSubdomainExists(next http.Handler) http.Handler
	AuthOrganizationIfSubdomainExists(next http.Handler) http.Handler
	SetSubdomain(next http.Handler) http.Handler
}

type MiddlewareCE struct {
	db.Repository
}

func NewMiddlewareCE(r db.Repository) *MiddlewareCE {
	return &MiddlewareCE{r}
}

// authenticateUser handles common user authentication logic and returns the authenticated user.
func (m *MiddlewareCE) authenticateUser(w http.ResponseWriter, r *http.Request) (*user.User, error) {
	ctx := r.Context()

	xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
	if xsrfTokenHeader == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token"))
	}

	xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	token, err := r.Cookie("access_token")
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	c, err := m.validateUserToken(token.Value)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	userID, err := uuid.FromString(c.UserID)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	u, err := m.User().Get(ctx, user.ByID(userID))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return u, nil
}

func (m *MiddlewareCE) getSubdomainIfCloudEdition(r *http.Request) (string, error) {
	if !config.Config.IsCloudEdition {
		return "", nil
	}
	return httpx.GetSubdomainFromHost(r.Host)
}

func (m *MiddlewareCE) validateOrganizationAccess(ctx context.Context, userID uuid.UUID, subdomain string) error {
	orgAccessOpts := []user.OrganizationAccessRepositoryOption{
		user.OrganizationAccessByUserID(userID),
	}
	if config.Config.IsCloudEdition {
		orgAccessOpts = append(orgAccessOpts, user.OrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := m.User().GetOrganizationAccess(ctx, orgAccessOpts...)
	return err
}

func (m *MiddlewareCE) getCurrentOrganization(ctx context.Context, subdomain string) (*organization.Organization, error) {
	opts := []organization.RepositoryOption{}
	if subdomain != "" && subdomain != "auth" {
		opts = append(opts, organization.BySubdomain(subdomain))
	}

	return m.Organization().Get(ctx, opts...)
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

func validateXSRFToken(header, cookie, claimToken string) error {
	if header == "" || cookie == "" || claimToken == "" {
		return errors.New("failed to get XSRF token")
	}
	if header != cookie && header != claimToken {
		return errors.New("invalid XSRF token")
	}
	return nil
}

func (m *MiddlewareCE) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := m.authenticateUser(w, r)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareCE) AuthUserWithOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := m.authenticateUser(w, r)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)

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

		if err := m.validateOrganizationAccess(ctx, u.ID, subdomain); err != nil {
			httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareCE) AuthUserWithOrganizationIfSubdomainExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := m.authenticateUser(w, r)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)

		subdomain, err := m.getSubdomainIfCloudEdition(r)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if subdomain != "" && subdomain != "auth" {
			o, err := m.getCurrentOrganization(ctx, subdomain)
			if err != nil {
				httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}

			if err := m.validateOrganizationAccess(ctx, u.ID, subdomain); err != nil {
				httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}

			ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareCE) AuthOrganizationIfSubdomainExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		subdomain, err := m.getSubdomainIfCloudEdition(r)
		if err != nil {
			httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if subdomain != "" && subdomain != "auth" {
			o, err := m.getCurrentOrganization(ctx, subdomain)
			if err != nil {
				httpx.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}

			ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareCE) SetSubdomain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomain, _ := httpx.GetSubdomainFromHost(r.Host)
		ctx := context.WithValue(r.Context(), ctxutil.SubdomainCtxKey, subdomain)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
