package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/dto/http/requests"
	"github.com/trysourcetool/sourcetool/backend/dto/http/responses"
	"github.com/trysourcetool/sourcetool/backend/organization/service"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type OrganizationHandler struct {
	service service.OrganizationService
}

func NewOrganizationHandler(service service.OrganizationService) *OrganizationHandler {
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

	out, err := h.service.Create(r.Context(), adapters.CreateOrganizationRequestToInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.CreateOrganizationOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// CheckSubdomainAvailability godoc
// @ID check-organization-subdomain-availability
// @Accept json
// @Produce json
// @Tags organizations
// @Param subdomain query string true "Subdomain to check for availability"
// @Success 200 {object} responses.StatusResponse
// @Failure 400 {object} errdefs.Error "Invalid subdomain format"
// @Failure 409 {object} errdefs.Error "Subdomain already exists"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /organizations/checkSubdomainAvailability [get].
func (h *OrganizationHandler) CheckSubdomainAvailability(w http.ResponseWriter, r *http.Request) {
	req := requests.CheckSubdomainAvailablityRequest{
		Subdomain: r.URL.Query().Get("subdomain"),
	}
	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	err := h.service.CheckSubdomainAvailability(r.Context(), adapters.CheckSubdomainAvailabilityRequestToInput(req))
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
