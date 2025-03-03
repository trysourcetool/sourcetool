package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
)

type GroupHandler struct {
	service group.Service
}

func NewGroupHandler(service group.Service) *GroupHandler {
	return &GroupHandler{service}
}

// Get godoc
// @ID get-group
// @Accept json
// @Produce json
// @Tags groups
// @Param groupID path string true "Group ID"
// @Success 200 {object} responses.GetGroupResponse
// @Failure default {object} errdefs.Error
// @Router /groups/{groupID} [get].
func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	req := requests.GetGroupRequest{
		GroupID: chi.URLParam(r, "groupID"),
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Get(r.Context(), adapters.GetGroupRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.GetGroupOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// List godoc
// @ID list-groups
// @Accept json
// @Produce json
// @Tags groups
// @Success 200 {object} responses.ListGroupsResponse
// @Failure default {object} errdefs.Error
// @Router /groups [get].
func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.ListGroupsOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Create godoc
// @ID create-group
// @Accept json
// @Produce json
// @Tags groups
// @Param Body body requests.CreateGroupRequest true " "
// @Success 200 {object} responses.CreateGroupResponse
// @Failure default {object} errdefs.Error
// @Router /groups [post].
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), adapters.CreateGroupRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.CreateGroupOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Update godoc
// @ID update-group
// @Accept json
// @Produce json
// @Tags groups
// @Param Body body requests.UpdateGroupRequest true " "
// @Param groupID path string true "Group ID"
// @Success 200 {object} responses.UpdateGroupResponse
// @Failure default {object} errdefs.Error
// @Router /groups/{groupID} [put].
func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateGroupRequest{
		GroupID: chi.URLParam(r, "groupID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), adapters.UpdateGroupRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.UpdateGroupOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Delete godoc
// @ID delete-group
// @Accept json
// @Produce json
// @Tags groups
// @Param groupID path string true "Group ID"
// @Success 200 {object} responses.DeleteGroupResponse
// @Failure default {object} errdefs.Error
// @Router /groups/{groupID} [delete].
func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	req := requests.DeleteGroupRequest{
		GroupID: chi.URLParam(r, "groupID"),
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Delete(r.Context(), adapters.DeleteGroupRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.DeleteGroupOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}
