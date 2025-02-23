package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type HostInstanceHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type HostInstanceHandlerCE struct {
	service hostinstance.Service
}

func NewHostInstanceHandlerCE(service hostinstance.Service) *HostInstanceHandlerCE {
	return &HostInstanceHandlerCE{service}
}

// Ping godoc
// @ID ping-host-instance
// @Accept json
// @Produce json
// @Tags hostInstances
// @Param pageId query string true "Page ID"
// @Success 200 {object} types.PingHostInstancePayload
// @Failure default {object} errdefs.Error
// @Router /hostInstances/ping [get].
func (h *HostInstanceHandlerCE) Ping(w http.ResponseWriter, r *http.Request) {
	in := types.PingHostInstanceInput{
		PageID: conv.NilValue(r.URL.Query().Get("pageId")),
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Ping(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
