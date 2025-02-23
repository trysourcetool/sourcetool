package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/organization"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type OrganizationHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	CheckSubdomainAvailability(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type OrganizationHandlerCE struct {
	service organization.Service
}

func NewOrganizationHandlerCE(service organization.Service) *OrganizationHandlerCE {
	return &OrganizationHandlerCE{service}
}

// Create godoc
// @ID create-organization
// @Accept json
// @Produce json
// @Tags organizations
// @Param Body body types.CreateOrganizationInput true " "
// @Success 200 {object} types.CreateOrganizationPayload
// @Failure default {object} errdefs.Error
// @Router /organizations [post].
func (h *OrganizationHandlerCE) Create(w http.ResponseWriter, r *http.Request) {
	var in types.CreateOrganizationInput
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

// CheckSubdomainAvairability godoc
// @ID check-organization-subdomain-availability
// @Accept json
// @Produce json
// @Tags organizations
// @Param subdomain query string true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /organizations/checkSubdomainAvailability [get].
func (h *OrganizationHandlerCE) CheckSubdomainAvailability(w http.ResponseWriter, r *http.Request) {
	in := types.CheckSubdomainAvailablityInput{
		Subdomain: r.URL.Query().Get("subdomain"),
	}
	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.CheckSubdomainAvailability(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdateUser godoc
// @ID update-organization-user
// @Accept json
// @Produce json
// @Tags organizations
// @Param userID path string true " "
// @Param Body body types.UpdateOrganizationUserInput true " "
// @Success 200 {object} types.UserPayload
// @Failure default {object} errdefs.Error
// @Router /organizations/users/{userID} [put].
func (h *OrganizationHandlerCE) UpdateUser(w http.ResponseWriter, r *http.Request) {
	in := types.UpdateOrganizationUserInput{
		UserID: chi.URLParam(r, "userID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateUser(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
