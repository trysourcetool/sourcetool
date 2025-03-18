package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
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
// @Success 200 {object} responses.ListPagesResponse
// @Failure default {object} errdefs.Error
// @Router /pages [get].
func (h *PageHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ListPagesOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
