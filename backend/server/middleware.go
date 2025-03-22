package server

import (
	"context"
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

func setSubdomain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomain, _ := httputil.GetSubdomainFromHost(r.Host)
		ctx := context.WithValue(r.Context(), ctxutil.SubdomainCtxKey, subdomain)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
