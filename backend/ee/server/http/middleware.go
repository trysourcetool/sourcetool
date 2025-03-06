package http

import (
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/server/http"
)

type MiddlewareEE struct {
	*http.MiddlewareCE
}

func NewMiddlewareEE(s infra.Store) *MiddlewareEE {
	return &MiddlewareEE{
		MiddlewareCE: http.NewMiddlewareCE(s),
	}
}
