package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type OrganizationHandler struct {
	service organization.Service
}

func NewOrganizationHandler(service organization.Service) *OrganizationHandler {
	return &OrganizationHandler{service}
}

// Create godoc
// @ID create-organization
// @Accept json
// @Produce json
// @Tags organizations
// @Param Body body requests.CreateOrganizationRequest true " "
// @Success 200 {object} responses.CreateOrganizationResponse
// @Failure default {object} errdefs.Error
// @Router /organizations [post].
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), adapters.CreateOrganizationRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.CreateOrganizationOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// CheckSubdomainAvairability godoc
// @ID check-organization-subdomain-availability
// @Accept json
// @Produce json
// @Tags organizations
// @Param subdomain query string true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /organizations/checkSubdomainAvailability [get].
func (h *OrganizationHandler) CheckSubdomainAvailability(w http.ResponseWriter, r *http.Request) {
	req := requests.CheckSubdomainAvailablityRequest{
		Subdomain: r.URL.Query().Get("subdomain"),
	}
	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	err := h.service.CheckSubdomainAvailability(r.Context(), adapters.CheckSubdomainAvailabilityRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	response := &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Subdomain is available",
	}

	if err := httputil.WriteJSON(w, http.StatusOK, response); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdateUser godoc
// @ID update-organization-user
// @Accept json
// @Produce json
// @Tags organizations
// @Param userID path string true " "
// @Param Body body requests.UpdateOrganizationUserRequest true " "
// @Success 200 {object} responses.UpdateUserResponse
// @Failure default {object} errdefs.Error
// @Router /organizations/users/{userID} [put].
func (h *OrganizationHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateOrganizationUserRequest{
		UserID: chi.URLParam(r, "userID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateUser(r.Context(), adapters.UpdateOrganizationUserRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateOrganizationUserOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
