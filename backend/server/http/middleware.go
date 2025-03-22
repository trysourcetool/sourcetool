package http

import (
	"context"
	"errors"
	"net/http"

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
	AuthUser(next http.Handler) http.Handler
	AuthUserWithOrganization(next http.Handler) http.Handler
	AuthUserWithOrganizationIfSubdomainExists(next http.Handler) http.Handler
	AuthOrganizationIfSubdomainExists(next http.Handler) http.Handler
	SetSubdomain(next http.Handler) http.Handler
}

type MiddlewareCE struct {
	infra.Store
}

func NewMiddlewareCE(s infra.Store) *MiddlewareCE {
	return &MiddlewareCE{s}
}

// authenticateUser handles common user authentication logic and returns the authenticated user.
func (m *MiddlewareCE) authenticateUser(w http.ResponseWriter, r *http.Request) (*model.User, error) {
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

	u, err := m.Store.User().Get(ctx, storeopts.UserByID(userID))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return u, nil
}

func (m *MiddlewareCE) getSubdomainIfCloudEdition(r *http.Request) (string, error) {
	if !config.Config.IsCloudEdition {
		return "", nil
	}
	return httputil.GetSubdomainFromHost(r.Host)
}

func (m *MiddlewareCE) validateOrganizationAccess(ctx context.Context, userID uuid.UUID, subdomain string) error {
	orgAccessOpts := []storeopts.UserOrganizationAccessOption{
		storeopts.UserOrganizationAccessByUserID(userID),
	}
	if config.Config.IsCloudEdition {
		orgAccessOpts = append(orgAccessOpts, storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := m.Store.User().GetOrganizationAccess(ctx, orgAccessOpts...)
	return err
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
			httputil.WriteErrJSON(ctx, w, err)
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
			httputil.WriteErrJSON(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)

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

		if err := m.validateOrganizationAccess(ctx, u.ID, subdomain); err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
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
			httputil.WriteErrJSON(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)

		subdomain, err := m.getSubdomainIfCloudEdition(r)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if subdomain != "" && subdomain != "auth" {
			o, err := m.getCurrentOrganization(ctx, subdomain)
			if err != nil {
				httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}

			if err := m.validateOrganizationAccess(ctx, u.ID, subdomain); err != nil {
				httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
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
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if subdomain != "" && subdomain != "auth" {
			o, err := m.getCurrentOrganization(ctx, subdomain)
			if err != nil {
				httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}

			ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareCE) SetSubdomain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomain, _ := httputil.GetSubdomainFromHost(r.Host)
		ctx := context.WithValue(r.Context(), ctxutil.SubdomainCtxKey, subdomain)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
