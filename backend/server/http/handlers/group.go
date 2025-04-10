package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
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

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Get(r.Context(), adapters.GetGroupRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.GetGroupOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
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
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ListGroupsOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Create godoc
// @ID create-group
// @Accept json
// @Produce json
// @Tags groups
// @Param Body body requests.CreateGroupRequest true "Group creation data including name and members"
// @Success 200 {object} responses.CreateGroupResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 403 {object} errdefs.Error "Insufficient permissions"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /groups [post].
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Create(r.Context(), adapters.CreateGroupRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.CreateGroupOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Update godoc
// @ID update-group
// @Accept json
// @Produce json
// @Tags groups
// @Param Body body requests.UpdateGroupRequest true "Group update data including name and members"
// @Param groupID path string true "Group ID to update"
// @Success 200 {object} responses.UpdateGroupResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 403 {object} errdefs.Error "Insufficient permissions"
// @Failure 404 {object} errdefs.Error "Group not found"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /groups/{groupID} [put].
func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateGroupRequest{
		GroupID: chi.URLParam(r, "groupID"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), adapters.UpdateGroupRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateGroupOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
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

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Delete(r.Context(), adapters.DeleteGroupRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.DeleteGroupOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
