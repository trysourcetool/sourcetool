package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/page"
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
// @Success 200 {object} types.ListPagesPayload
// @Failure default {object} errdefs.Error
// @Router /pages [get].
func (h *PageHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
