package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authn"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
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
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token")))
			return
		}

		xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		token, err := r.Cookie("access_token")
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		c, err := m.validateUserToken(token.Value)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		userID, err := uuid.FromString(c.UserID)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		u, err := m.Store.User().Get(ctx, storeopts.UserByID(userID))
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutils.CurrentUserCtxKey, u)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCE) AuthOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		subdomain := strings.Split(r.Host, ".")[0]
		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutils.CurrentOrganizationCtxKey, o)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCE) AuthUserWithOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
		if xsrfTokenHeader == "" {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token")))
			return
		}

		xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		token, err := r.Cookie("access_token")
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		c, err := m.validateUserToken(token.Value)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		userID, err := uuid.FromString(c.UserID)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		u, err := m.Store.User().Get(ctx, storeopts.UserByID(userID))
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutils.CurrentUserCtxKey, u)
		r = r.WithContext(ctx)

		subdomain := strings.Split(r.Host, ".")[0]
		o, err := m.getCurrentOrganization(ctx, subdomain)
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		_, err = m.Store.User().GetOrganizationAccess(ctx, storeopts.UserOrganizationAccessByUserID(u.ID), storeopts.UserOrganizationAccessByOrganizationSubdomain(subdomain))
		if err != nil {
			httputils.WriteErrJSON(ctx, w, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, ctxutils.CurrentOrganizationCtxKey, o)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareCE) getCurrentOrganization(ctx context.Context, subdomain string) (*model.Organization, error) {
	o, err := m.Store.Organization().Get(ctx, storeopts.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (m *MiddlewareCE) validateUserToken(token string) (*authn.UserAuthClaims, error) {
	if token == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get token"))
	}

	claims, err := authn.ParseToken[*authn.UserAuthClaims](token)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return claims, nil
}

func (m *MiddlewareCE) SetHTTPHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		ctx := context.WithValue(r.Context(), ctxutils.HTTPHostCtxKey, host)
		r = r.WithContext(ctx)

		referer := r.Header["Referer"]
		if len(referer) != 0 {
			ctx := context.WithValue(r.Context(), ctxutils.HTTPRefererCtxKey, referer[0])
			r = r.WithContext(ctx)
		}
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
