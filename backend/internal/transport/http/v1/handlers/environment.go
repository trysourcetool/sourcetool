package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/mapper"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/pkg/httpx"
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

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Get(r.Context(), mapper.GetEnvironmentRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.GetEnvironmentOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.ListEnvironmentsOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Create godoc
// @ID create-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param Body body requests.CreateEnvironmentRequest true "Environment creation data including name and configuration"
// @Success 200 {object} responses.CreateEnvironmentResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 403 {object} errdefs.Error "Insufficient permissions"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /environments [post].
func (h *EnvironmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateEnvironmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), mapper.CreateEnvironmentRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.CreateEnvironmentOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Update godoc
// @ID update-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param Body body requests.UpdateEnvironmentRequest true "Environment update data including name and configuration"
// @Param environmentID path string true "Environment ID to update"
// @Success 200 {object} responses.UpdateEnvironmentResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 403 {object} errdefs.Error "Insufficient permissions"
// @Failure 404 {object} errdefs.Error "Environment not found"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /environments/{environmentID} [put].
func (h *EnvironmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateEnvironmentRequest{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), mapper.UpdateEnvironmentRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.UpdateEnvironmentOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Delete(r.Context(), mapper.DeleteEnvironmentRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.DeleteEnvironmentOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}
}
