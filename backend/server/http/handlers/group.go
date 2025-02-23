package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type GroupHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type GroupHandlerCE struct {
	service group.Service
}

func NewGroupHandlerCE(service group.Service) *GroupHandlerCE {
	return &GroupHandlerCE{service}
}

// Get godoc
// @ID get-group
// @Accept json
// @Produce json
// @Tags groups
// @Param groupID path string true "Group ID"
// @Success 200 {object} types.GetGroupPayload
// @Failure default {object} errdefs.Error
// @Router /groups/{groupID} [get].
func (h *GroupHandlerCE) Get(w http.ResponseWriter, r *http.Request) {
	in := types.GetGroupInput{
		GroupID: chi.URLParam(r, "groupID"),
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
// @ID list-groups
// @Accept json
// @Produce json
// @Tags groups
// @Success 200 {object} types.ListGroupsPayload
// @Failure default {object} errdefs.Error
// @Router /groups [get].
func (h *GroupHandlerCE) List(w http.ResponseWriter, r *http.Request) {
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
// @ID create-group
// @Accept json
// @Produce json
// @Tags groups
// @Param Body body types.CreateGroupInput true " "
// @Success 200 {object} types.CreateGroupPayload
// @Failure default {object} errdefs.Error
// @Router /groups [post].
func (h *GroupHandlerCE) Create(w http.ResponseWriter, r *http.Request) {
	var in types.CreateGroupInput
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
// @ID update-group
// @Accept json
// @Produce json
// @Tags groups
// @Param Body body types.UpdateGroupInput true " "
// @Param groupID path string true "Group ID"
// @Success 200 {object} types.UpdateGroupPayload
// @Failure default {object} errdefs.Error
// @Router /groups/{groupID} [put].
func (h *GroupHandlerCE) Update(w http.ResponseWriter, r *http.Request) {
	in := types.UpdateGroupInput{
		GroupID: chi.URLParam(r, "groupID"),
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
// @ID delete-group
// @Accept json
// @Produce json
// @Tags groups
// @Param groupID path string true "Group ID"
// @Success 200 {object} types.DeleteGroupPayload
// @Failure default {object} errdefs.Error
// @Router /groups/{groupID} [delete].
func (h *GroupHandlerCE) Delete(w http.ResponseWriter, r *http.Request) {
	in := types.DeleteGroupInput{
		GroupID: chi.URLParam(r, "groupID"),
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
