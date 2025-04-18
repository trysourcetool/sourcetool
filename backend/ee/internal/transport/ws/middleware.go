package ws

import (
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/ws"
)

type MiddlewareEE struct {
	*ws.MiddlewareCE
}

func NewMiddlewareEE(r db.Repository) *MiddlewareEE {
	return &MiddlewareEE{
		MiddlewareCE: ws.NewMiddlewareCE(r),
	}
}
