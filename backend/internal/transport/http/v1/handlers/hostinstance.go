package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/app/hostinstance"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/render"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/mapper"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
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
		PageID: internal.NilValue(r.URL.Query().Get("pageId")),
	}

	if err := validateRequest(req); err != nil {
		render.Error(r.Context(), w, err)
		return
	}

	out, err := h.service.Ping(r.Context(), mapper.PingHostInstanceRequestToInput(req))
	if err != nil {
		render.Error(r.Context(), w, err)
		return
	}

	if err := render.JSON(w, http.StatusOK, mapper.PingHostInstanceOutputToResponse(out)); err != nil {
		render.Error(r.Context(), w, err)
		return
	}
}
