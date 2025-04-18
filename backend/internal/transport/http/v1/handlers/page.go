package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/internal/app/page"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/mapper"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/pkg/httpx"
)

type PageHandler struct {
	service page.Service
}

func NewPageHandler(service page.Service) *PageHandler {
	return &PageHandler{service}
}

// List godoc
// @ID list-pages
// @Accept json
// @Produce json
// @Tags pages
// @Param environmentId query string true "Environment ID"
// @Success 200 {object} responses.ListPagesResponse
// @Failure default {object} errdefs.Error
// @Router /pages [get].
func (h *PageHandler) List(w http.ResponseWriter, r *http.Request) {
	in := requests.ListPagesRequest{
		EnvironmentID: r.URL.Query().Get("environmentId"),
	}
	if err := httpx.ValidateRequest(in); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.List(r.Context(), mapper.ListPagesRequestToInput(in))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.ListPagesOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}
}
