package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
)

type HostInstanceHandler struct {
	service hostinstance.Service
}

func NewHostInstanceHandler(service hostinstance.Service) *HostInstanceHandler {
	return &HostInstanceHandler{service}
}

// Ping godoc
// @ID ping-host-instance
// @Accept json
// @Produce json
// @Tags hostInstances
// @Param pageId query string true "Page ID"
// @Success 200 {object} responses.PingHostInstanceResponse
// @Failure default {object} errdefs.Error
// @Router /hostInstances/ping [get].
func (h *HostInstanceHandler) Ping(w http.ResponseWriter, r *http.Request) {
	req := requests.PingHostInstanceRequest{
		PageID: conv.NilValue(r.URL.Query().Get("pageId")),
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Ping(r.Context(), adapters.PingHostInstanceRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.PingHostInstanceOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
