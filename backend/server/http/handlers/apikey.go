package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/apikey/service"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type APIKeyHandler struct {
	service service.APIKeyService
}

func NewAPIKeyHandler(service service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service}
}

// Get godoc
// @ID get-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} responses.GetAPIKeyResponse
// @Failure default {object} errdefs.Error
// @Router /apiKeys/{apiKeyID} [get].
func (h *APIKeyHandler) Get(w http.ResponseWriter, r *http.Request) {
	req := requests.GetAPIKeyRequest{
		APIKeyID: chi.URLParam(r, "apiKeyID"),
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Get(r.Context(), adapters.GetAPIKeyRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.GetAPIKeyOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// List godoc
// @ID list-apikeys
// @Accept json
// @Produce json
// @Tags apiKeys
// @Success 200 {object} responses.ListAPIKeysResponse
// @Failure default {object} errdefs.Error
// @Router /apiKeys [get].
func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ListAPIKeysOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Create godoc
// @ID create-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param Body body requests.CreateAPIKeyRequest true "API key creation data including name and expiration"
// @Success 200 {object} responses.CreateAPIKeyResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 403 {object} errdefs.Error "Insufficient permissions"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /apiKeys [post].
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), adapters.CreateAPIKeyRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.CreateAPIKeyOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Update godoc
// @ID update-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param Body body requests.UpdateAPIKeyRequest true "API key update data including name and status"
// @Param apiKeyID path string true "API Key ID to update"
// @Success 200 {object} responses.UpdateAPIKeyResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 403 {object} errdefs.Error "Insufficient permissions"
// @Failure 404 {object} errdefs.Error "API key not found"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /apiKeys/{apiKeyID} [put].
func (h *APIKeyHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateAPIKeyRequest{
		APIKeyID: chi.URLParam(r, "apiKeyID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), adapters.UpdateAPIKeyRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateAPIKeyOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Delete godoc
// @ID delete-apikey
// @Accept json
// @Produce json
// @Tags apiKeys
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} responses.DeleteAPIKeyResponse
// @Failure default {object} errdefs.Error
// @Router /apiKeys/{apiKeyID} [delete].
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	req := requests.DeleteAPIKeyRequest{
		APIKeyID: chi.URLParam(r, "apiKeyID"),
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Delete(r.Context(), adapters.DeleteAPIKeyRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.DeleteAPIKeyOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
