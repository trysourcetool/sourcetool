package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/hostinstance/service"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/utils/conv"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type HostInstanceHandler struct {
	service service.HostInstanceService
}

func NewHostInstanceHandler(service service.HostInstanceService) *HostInstanceHandler {
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

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Ping(r.Context(), adapters.PingHostInstanceRequestToInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.PingHostInstanceOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
