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
	AuthOrganization(next http.Handler) http.Handler
	AuthUserWithOrganization(next http.Handler) http.Handler
	SetHTTPHeader(next http.Handler) http.Handler
}

type MiddlewareCE struct {
	infra.Store
}

func NewMiddlewareCE(s infra.Store) *MiddlewareCE {
	return &MiddlewareCE{s}
}

func (m *MiddlewareCE) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
		if xsrfTokenHeader == "" {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token")))
			return
		}

		xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		token, err := r.Cookie("access_token")
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		c, err := m.validateUserToken(token.Value)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		userID, err := uuid.FromString(c.UserID)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		u, err := m.Store.User().Get(ctx, storeopts.UserByID(userID))
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCE) AuthOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var subdomain string
		var err error
		if config.Config.IsCloudEdition {
			subdomain, err = httputil.GetSubdomainFromHost(r.Host)
			if err != nil {
				httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}
		}

		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCE) AuthUserWithOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
		if xsrfTokenHeader == "" {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token")))
			return
		}

		xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		token, err := r.Cookie("access_token")
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		c, err := m.validateUserToken(token.Value)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		userID, err := uuid.FromString(c.UserID)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		u, err := m.Store.User().Get(ctx, storeopts.UserByID(userID))
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentUserCtxKey, u)
		r = r.WithContext(ctx)

		var subdomain string
		if config.Config.IsCloudEdition {
			subdomain, err = httputil.GetSubdomainFromHost(r.Host)
			if err != nil {
				httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
				return
			}
		}

		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		_, err = m.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(u.ID), storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
		if err != nil {
			httputil.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutil.CurrentOrganizationCtxKey, o)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCE) getCurrentOrganization(ctx context.Context, subdomain string) (*model.Organization, error) {
	opts := []storeopts.OrganizationOption{}
	if subdomain != "" {
		opts = append(opts, storeopts.OrganizationBySubdomain(subdomain))
	}

	o, err := m.Store.Organization().Get(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return o, nil
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

func (m *MiddlewareCE) SetHTTPHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		ctx := context.WithValue(r.Context(), ctxutil.HTTPHostCtxKey, host)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func validateXSRFToken(a, b, c string) error {
	if a == "" || b == "" || c == "" {
		return errors.New("failed to get XSRF token")
	}
	if a != b && a != c {
		return errors.New("invalid XSRF token")
	}
	return nil
}
