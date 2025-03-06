package ws

import (
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/server/ws"
)

type MiddlewareEE struct {
	*ws.MiddlewareCE
}

func NewMiddlewareEE(s infra.Store) *MiddlewareEE {
	return &MiddlewareEE{
		MiddlewareCE: ws.NewMiddlewareCE(s),
	}
}
