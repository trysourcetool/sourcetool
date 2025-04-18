package v1

import (
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	v1 "github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1"
)

type MiddlewareEE struct {
	*v1.MiddlewareCE
}

func NewMiddlewareEE(d *infra.Dependency) *MiddlewareEE {
	return &MiddlewareEE{
		MiddlewareCE: v1.NewMiddlewareCE(d.Repository),
	}
}
