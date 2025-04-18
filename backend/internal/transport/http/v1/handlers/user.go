package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/user"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/mapper"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/requests"
	"github.com/trysourcetool/sourcetool/backend/internal/transport/http/v1/responses"
	"github.com/trysourcetool/sourcetool/backend/pkg/httpx"
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.GetMeOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateMe(r.Context(), mapper.UpdateMeRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.UpdateMeOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.SendUpdateMeEmailInstructions(r.Context(), mapper.SendUpdateMeEmailInstructionsRequestToInput(req)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully sent update email instructions",
	}); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateMeEmail(r.Context(), mapper.UpdateMeEmailRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.UpdateMeEmailOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.ListUsersOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), mapper.UpdateUserRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.UpdateUserOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.Delete(r.Context(), mapper.DeleteUserRequestToInput(req)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	response := &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully deleted user",
	}

	if err := httpx.WriteJSON(w, http.StatusOK, response); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// ResendUserInvitation godoc
// @ID resend-user-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param invitationID path string true " "
// @Success 200 {object} responses.ResendUserInvitationResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/{invitationID}/resend [post].
func (h *UserHandler) ResendUserInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.ResendUserInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.ResendUserInvitation(r.Context(), mapper.ResendUserInvitationRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.ResendUserInvitationOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
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
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.ValidateRequest(req); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.CreateUserInvitations(r.Context(), mapper.CreateUserInvitationsRequestToInput(req))
	if err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, mapper.CreateUserInvitationsOutputToResponse(out)); err != nil {
		httpx.WriteErrJSON(r.Context(), w, err)
		return
	}
}
