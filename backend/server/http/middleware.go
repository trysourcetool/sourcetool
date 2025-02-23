package http

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

		u, err := m.Store.User().Get(ctx, model.UserByEmail(c.Email))
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

		u, err := m.Store.User().Get(ctx, model.UserByEmail(c.Email))
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

		_, err = m.Store.User().GetOrganizationAccess(ctx, model.UserOrganizationAccessByUserID(u.ID), model.UserOrganizationAccessByOrganizationSubdomain(subdomain))
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
	o, err := m.Store.Organization().Get(ctx, model.OrganizationBySubdomain(subdomain))
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (m *MiddlewareCE) validateUserToken(token string) (*model.UserClaims, error) {
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
