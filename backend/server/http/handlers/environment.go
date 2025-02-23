package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
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
// @Success 200 {object} types.GetEnvironmentPayload
// @Failure default {object} errdefs.Error
// @Router /environments/{environmentID} [get].
func (h *EnvironmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	in := types.GetEnvironmentInput{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Get(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// List godoc
// @ID list-environments
// @Accept json
// @Produce json
// @Tags environments
// @Success 200 {object} types.ListEnvironmentsPayload
// @Failure default {object} errdefs.Error
// @Router /environments [get].
func (h *EnvironmentHandler) List(w http.ResponseWriter, r *http.Request) {
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

// Create godoc
// @ID create-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param Body body types.CreateEnvironmentInput true " "
// @Success 200 {object} types.CreateEnvironmentPayload
// @Failure default {object} errdefs.Error
// @Router /environments [post].
func (h *EnvironmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in types.CreateEnvironmentInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Update godoc
// @ID update-environment
// @Accept json
// @Produce json
// @Tags environments
// @Param Body body types.UpdateEnvironmentInput true " "
// @Param environmentID path string true "Environment ID"
// @Success 200 {object} types.UpdateEnvironmentPayload
// @Failure default {object} errdefs.Error
// @Router /environments/{environmentID} [put].
func (h *EnvironmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	in := types.UpdateEnvironmentInput{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
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
// @Success 200 {object} types.DeleteEnvironmentPayload
// @Failure default {object} errdefs.Error
// @Router /environments/{environmentID} [delete].
func (h *EnvironmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	in := types.DeleteEnvironmentInput{
		EnvironmentID: chi.URLParam(r, "environmentID"),
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Delete(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
