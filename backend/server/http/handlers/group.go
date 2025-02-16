package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type GroupHandlerCE interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	// AddUsers(w http.ResponseWriter, r *http.Request)
	// RemoveUser(w http.ResponseWriter, r *http.Request)
}

type GroupHandlerCEImpl struct {
	service group.ServiceCE
}

func NewGroupHandlerCE(service group.ServiceCE) *GroupHandlerCEImpl {
	return &GroupHandlerCEImpl{service}
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
func (h *GroupHandlerCEImpl) Get(w http.ResponseWriter, r *http.Request) {
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
func (h *GroupHandlerCEImpl) List(w http.ResponseWriter, r *http.Request) {
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
func (h *GroupHandlerCEImpl) Create(w http.ResponseWriter, r *http.Request) {
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
func (h *GroupHandlerCEImpl) Update(w http.ResponseWriter, r *http.Request) {
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
func (h *GroupHandlerCEImpl) Delete(w http.ResponseWriter, r *http.Request) {
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

// TODO: Implement later when UI is ready
// // AddUsers godoc
// // @ID add-users-to-group
// // @Accept json
// // @Produce json
// // @Tags groups
// // @Param groupID path string true "Group ID"
// // @Param Body body types.AddUsersToGroupInput true " "
// // @Success 200 {object} types.AddUsersToGroupPayload
// // @Failure default {object} errdefs.Error
// // @Router /groups/{groupID}/users [post].
// func (h *GroupHandlerCEImpl) AddUsers(w http.ResponseWriter, r *http.Request) {
// 	in := types.AddUsersToGroupInput{
// 		GroupID: chi.URLParam(r, "groupID"),
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
//
// 	if err := httputils.ValidateRequest(in); err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
//
// 	out, err := h.service.AddUsers(r.Context(), in)
// 	if err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
//
// 	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
// }
//
// // RemoveUser godoc
// // @ID remove-user-from-group
// // @Accept json
// // @Produce json
// // @Tags groups
// // @Param groupID path string true "Group ID"
// // @Param userID path string true "User ID"
// // @Success 200 {object} types.RemoveUserFromGroupInput
// // @Failure default {object} errdefs.Error
// // @Router /groups/{groupID}/users/{userID} [delete].
// func (h *GroupHandlerCEImpl) RemoveUser(w http.ResponseWriter, r *http.Request) {
// 	in := types.RemoveUserFromGroupInput{
// 		GroupID: chi.URLParam(r, "groupID"),
// 		UserID:  chi.URLParam(r, "userID"),
// 	}
//
// 	if err := httputils.ValidateRequest(in); err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
//
// 	out, err := h.service.RemoveUser(r.Context(), in)
// 	if err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
//
// 	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
// 		httputils.WriteErrJSON(r.Context(), w, err)
// 		return
// 	}
// }
//
