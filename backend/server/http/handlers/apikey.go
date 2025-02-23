package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/apikey"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type APIKeyHandler struct {
	service apikey.Service
}

func NewAPIKeyHandler(service apikey.Service) *APIKeyHandler {
	return &APIKeyHandler{service}
}

// Get godoc
// @ID get-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} types.GetAPIKeyPayload
// @Failure default {object} errdefs.Error
// @Router /apiKeys/{apiKeyID} [get].
func (h *APIKeyHandler) Get(w http.ResponseWriter, r *http.Request) {
	in := types.GetAPIKeyInput{
		APIKeyID: chi.URLParam(r, "apiKeyID"),
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
// @ID list-apikeys
// @Accept json
// @Produce json
// @Tags apiKeys
// @Success 200 {object} types.ListAPIKeysPayload
// @Failure default {object} errdefs.Error
// @Router /apiKeys [get].
func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
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
// @ID create-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param Body body types.CreateAPIKeyInput true " "
// @Success 200 {object} types.CreateAPIKeyPayload
// @Failure default {object} errdefs.Error
// @Router /apiKeys [post].
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in types.CreateAPIKeyInput
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
// @ID update-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param Body body types.UpdateAPIKeyInput true " "
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} types.UpdateAPIKeyPayload
// @Failure default {object} errdefs.Error
// @Router /apiKeys/{apiKeyID} [put].
func (h *APIKeyHandler) Update(w http.ResponseWriter, r *http.Request) {
	in := types.UpdateAPIKeyInput{
		APIKeyID: chi.URLParam(r, "apiKeyID"),
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
// @ID delete-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} types.DeleteAPIKeyPayload
// @Failure default {object} errdefs.Error
// @Router /apiKeys/{apiKeyID} [delete].
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	in := types.DeleteAPIKeyInput{
		APIKeyID: chi.URLParam(r, "apiKeyID"),
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
