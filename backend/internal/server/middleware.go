package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/jwt"
)

// authenticateUser handles common user authentication logic and returns the authenticated user.
func (s *Server) authenticateUser(w http.ResponseWriter, r *http.Request) (*core.User, error) {
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

	c, err := s.validateUserToken(token.Value)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	userID, err := uuid.FromString(c.Subject)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	u, err := s.db.User().Get(ctx, database.UserByID(userID))
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return u, nil
}

func (s *Server) getSubdomainIfCloudEdition(r *http.Request) (string, error) {
	if !config.Config.IsCloudEdition {
		return "", nil
	}
	return internal.GetSubdomainFromHost(r.Host)
}

func (s *Server) validateOrganizationAccess(ctx context.Context, userID uuid.UUID, subdomain string) error {
	orgAccessOpts := []database.UserOrganizationAccessQuery{
		database.UserOrganizationAccessByUserID(userID),
	}
	if config.Config.IsCloudEdition {
		orgAccessOpts = append(orgAccessOpts, database.UserOrganizationAccessByOrganizationSubdomain(subdomain))
	}
	_, err := s.db.User().GetOrganizationAccess(ctx, orgAccessOpts...)
	return err
}

func (s *Server) getContextOrganization(ctx context.Context, subdomain string) (*core.Organization, error) {
	opts := []database.OrganizationQuery{}
	if subdomain != "" && subdomain != "auth" {
		opts = append(opts, database.OrganizationBySubdomain(subdomain))
	}

	return s.db.Organization().Get(ctx, opts...)
}

func (s *Server) validateUserToken(token string) (*jwt.AuthClaims, error) {
	if token == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get token"))
	}

	claims, err := jwt.ParseAuthClaims(token)
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

func (s *Server) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := s.authenticateUser(w, r)
		if err != nil {
			s.serveError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, internal.ContextUserKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) authUserWithOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := s.authenticateUser(w, r)
		if err != nil {
			s.serveError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, internal.ContextUserKey, u)

		subdomain, err := s.getSubdomainIfCloudEdition(r)
		if err != nil {
			s.serveError(w, r, errdefs.ErrUnauthenticated(err))
			return
		}

		o, err := s.getContextOrganization(ctx, subdomain)
		if err != nil {
			s.serveError(w, r, errdefs.ErrUnauthenticated(err))
			return
		}

		if err := s.validateOrganizationAccess(ctx, u.ID, subdomain); err != nil {
			s.serveError(w, r, errdefs.ErrUnauthenticated(err))
			return
		}

		ctx = context.WithValue(ctx, internal.ContextOrganizationKey, o)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) authUserWithOrganizationIfSubdomainExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := s.authenticateUser(w, r)
		if err != nil {
			s.serveError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, internal.ContextUserKey, u)

		subdomain, err := s.getSubdomainIfCloudEdition(r)
		if err != nil {
			s.serveError(w, r, errdefs.ErrUnauthenticated(err))
			return
		}

		if subdomain != "" && subdomain != "auth" {
			o, err := s.getContextOrganization(ctx, subdomain)
			if err != nil {
				s.serveError(w, r, errdefs.ErrUnauthenticated(err))
				return
			}

			if err := s.validateOrganizationAccess(ctx, u.ID, subdomain); err != nil {
				s.serveError(w, r, errdefs.ErrUnauthenticated(err))
				return
			}

			ctx = context.WithValue(ctx, internal.ContextOrganizationKey, o)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) authOrganizationIfSubdomainExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		subdomain, err := s.getSubdomainIfCloudEdition(r)
		if err != nil {
			s.serveError(w, r, errdefs.ErrUnauthenticated(err))
			return
		}

		if subdomain != "" && subdomain != "auth" {
			o, err := s.getContextOrganization(ctx, subdomain)
			if err != nil {
				s.serveError(w, r, errdefs.ErrUnauthenticated(err))
				return
			}

			ctx = context.WithValue(ctx, internal.ContextOrganizationKey, o)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) getCurrentUser(ctx context.Context, token string) (*core.User, error) {
	c, err := s.validateUserToken(token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(c.Subject)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	u, err := s.db.User().Get(ctx, database.UserByID(userID))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Server) setSubdomain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomain, _ := internal.GetSubdomainFromHost(r.Host)
		ctx := context.WithValue(r.Context(), internal.ContextSubdomainKey, subdomain)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) extractIncomingToken(headerValue string) (string, error) {
	if !strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
		return "", fmt.Errorf("invalid or malformed %q header, expected 'Bearer JWT-token...'", headerValue)
	}
	return strings.Split(headerValue, " ")[1], nil
}

func (s *Server) authWebSocketUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		subdomain, err := s.getSubdomainIfCloudEdition(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		o, err := s.getContextOrganization(ctx, subdomain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, internal.ContextOrganizationKey, o)

		if token, err := r.Cookie("access_token"); err == nil {
			u, err := s.getCurrentUser(ctx, token.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, internal.ContextUserKey, u)
		} else if apiKeyHeader := r.Header.Get("Authorization"); apiKeyHeader != "" {
			apikeyVal, err := s.extractIncomingToken(apiKeyHeader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			hashedAPIKey := core.HashAPIKey(apikeyVal)
			apikey, err := s.db.APIKey().Get(ctx, database.APIKeyByKeyHash(hashedAPIKey))
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
