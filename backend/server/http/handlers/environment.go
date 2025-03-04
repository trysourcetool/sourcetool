package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
)

type EnvironmentHandler struct {
	service environment.Service
}

func NewEnvironmentHandler(service environment.Service) *EnvironmentHandler {
	return &EnvironmentHandler{service}
}

// Get godoc
// @ID get-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param environmentID path string true "Environment ID"
// @Success 200 {object} responses.GetEnvironmentResponse
// @Failure default {object} errdefs.Error
// @Router /environments/{environmentID} [get].
func (h *EnvironmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	req := requests.GetEnvironmentRequest{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Get(r.Context(), adapters.GetEnvironmentRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.GetEnvironmentOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// List godoc
// @ID list-environments
// @Accept json
// @Produce json
// @Tags environments
// @Success 200 {object} responses.ListEnvironmentsResponse
// @Failure default {object} errdefs.Error
// @Router /environments [get].
func (h *EnvironmentHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.ListEnvironmentsOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Create godoc
// @ID create-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param Body body requests.CreateEnvironmentRequest true " "
// @Success 200 {object} responses.CreateEnvironmentResponse
// @Failure default {object} errdefs.Error
// @Router /environments [post].
func (h *EnvironmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), adapters.CreateEnvironmentRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.CreateEnvironmentOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Update godoc
// @ID update-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param Body body requests.UpdateEnvironmentRequest true " "
// @Param environmentID path string true "Environment ID"
// @Success 200 {object} responses.UpdateEnvironmentResponse
// @Failure default {object} errdefs.Error
// @Router /environments/{environmentID} [put].
func (h *EnvironmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateEnvironmentRequest{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), adapters.UpdateEnvironmentRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.UpdateEnvironmentOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Delete godoc
// @ID delete-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param environmentID path string true "Environment ID"
// @Success 200 {object} responses.DeleteEnvironmentResponse
// @Failure default {object} errdefs.Error
// @Router /environments/{environmentID} [delete].
func (h *EnvironmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	req := requests.DeleteEnvironmentRequest{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Delete(r.Context(), adapters.DeleteEnvironmentRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.DeleteEnvironmentOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
