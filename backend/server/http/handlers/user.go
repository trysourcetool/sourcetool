package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetMe godoc
// @ID get-me
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.GetMeResponse
// @Failure default {object} errdefs.Error
// @Router /users/me [get].
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.GetMe(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.GetMeOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdateMe godoc
// @ID update-me
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.UpdateMeRequest true " "
// @Success 200 {object} responses.UpdateMeResponse
// @Failure default {object} errdefs.Error
// @Router /users/me [put].
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateMe(r.Context(), adapters.UpdateMeRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateMeOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SendUpdateMeEmailInstructions godoc
// @ID send-update-me-email-instructions
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SendUpdateMeEmailInstructionsRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/me/email/instructions [post].
func (h *UserHandler) SendUpdateMeEmailInstructions(w http.ResponseWriter, r *http.Request) {
	var req requests.SendUpdateMeEmailInstructionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.SendUpdateMeEmailInstructions(r.Context(), adapters.SendUpdateMeEmailInstructionsRequestToDTOInput(req)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully sent update email instructions",
	}); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdateMeEmail godoc
// @ID update-me-email
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.UpdateMeEmailRequest true " "
// @Success 200 {object} responses.UpdateMeEmailResponse
// @Failure default {object} errdefs.Error
// @Router /users/me/email [put].
func (h *UserHandler) UpdateMeEmail(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateMeEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateMeEmail(r.Context(), adapters.UpdateMeEmailRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateMeEmailOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// List godoc
// @ID list-users
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.ListUsersResponse
// @Failure default {object} errdefs.Error
// @Router /users [get].
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ListUsersOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdateUser godoc
// @ID update-user
// @Accept json
// @Produce json
// @Tags users
// @Param userID path string true " "
// @Param Body body requests.UpdateUserRequest true " "
// @Success 200 {object} responses.UpdateUserResponse
// @Failure default {object} errdefs.Error
// @Router /users/{userID} [put].
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := requests.UpdateUserRequest{
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

	out, err := h.service.Update(r.Context(), adapters.UpdateUserRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateUserOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// DeleteUser godoc
// @ID delete-user
// @Accept json
// @Produce json
// @Tags users
// @Param userID path string true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/{userID} [delete].
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	req := requests.DeleteUserRequest{
		UserID: chi.URLParam(r, "userID"),
	}
	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.Delete(r.Context(), adapters.DeleteUserRequestToDTOInput(req)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	response := &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully deleted user",
	}

	if err := httputil.WriteJSON(w, http.StatusOK, response); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// ResendInvitation godoc
// @ID resend-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.ResendInvitationRequest true " "
// @Success 200 {object} responses.ResendInvitationResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/resend [post].
func (h *UserHandler) ResendUserInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.ResendUserInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.ResendUserInvitation(r.Context(), adapters.ResendUserInvitationRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ResendUserInvitationOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// CreateUserInvitations godoc
// @ID create-user-invitations
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.CreateUserInvitationsRequest true " "
// @Success 200 {object} responses.CreateUserInvitationsResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations [post].
func (h *UserHandler) CreateUserInvitations(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateUserInvitationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.CreateUserInvitations(r.Context(), adapters.CreateUserInvitationsRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.CreateUserInvitationsOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
