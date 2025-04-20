package ws

import (
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws"
)

type MiddlewareEE struct {
	*ws.MiddlewareCE
}

func NewMiddlewareEE(r port.Repository) *MiddlewareEE {
	return &MiddlewareEE{
		MiddlewareCE: ws.NewMiddlewareCE(r),
	}
}
