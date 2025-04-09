package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type UserHandler struct {
	service      user.Service
	cookieConfig *CookieConfig
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{
		service:      service,
		cookieConfig: NewCookieConfig(),
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

// Update godoc
// @ID update-user
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.UpdateUserRequest true " "
// @Success 200 {object} responses.UpdateUserResponse
// @Failure default {object} errdefs.Error
// @Router /users [put].
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateUserRequest
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

// SendUpdateEmailInstructions godoc
// @ID send-update-email-instructions
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SendUpdateUserEmailInstructionsRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/sendUpdateEmailInstructions [post].
func (h *UserHandler) SendUpdateEmailInstructions(w http.ResponseWriter, r *http.Request) {
	var req requests.SendUpdateUserEmailInstructionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.SendUpdateEmailInstructions(r.Context(), adapters.SendUpdateUserEmailInstructionsRequestToDTOInput(req)); err != nil {
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

// UpdateEmail godoc
// @ID update-user-email
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.UpdateUserEmailRequest true " "
// @Success 200 {object} responses.UpdateUserEmailResponse
// @Failure default {object} errdefs.Error
// @Router /users/email [put].
func (h *UserHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateUserEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateEmail(r.Context(), adapters.UpdateUserEmailRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.UpdateUserEmailOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RequestMagicLink godoc
// @ID request-magic-link
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.RequestMagicLinkRequest true " "
// @Success 200 {object} responses.RequestMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /users/signin/magic/request [post].
func (h *UserHandler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	res, err := h.service.RequestMagicLink(r.Context(), adapters.RequestMagicLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RequestMagicLinkOutputToResponse(res)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// AuthenticateWithMagicLink godoc
// @ID authenticate-with-magic-link
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.AuthenticateWithMagicLinkRequest true " "
// @Success 200 {object} responses.AuthenticateWithMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /users/auth/magic/authenticate [post].
func (h *UserHandler) AuthenticateWithMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.AuthenticateWithMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.AuthenticateWithMagicLink(r.Context(), adapters.AuthenticateWithMagicLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if !out.IsOrganizationExists {
		h.cookieConfig.SetTmpAuthCookie(w, out.Token, out.XSRFToken, config.Config.AuthDomain())
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.AuthenticateWithMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// @Router /users/auth/magic/register [post].
func (h *UserHandler) RegisterWithMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterWithMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrInvalidArgument(err))
		return
	}

	out, err := h.service.RegisterWithMagicLink(r.Context(), adapters.RegisterWithMagicLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if config.Config.IsCloudEdition {
		h.cookieConfig.SetTmpAuthCookie(w, out.Token, out.XSRFToken, config.Config.AuthDomain())
	} else {
		h.cookieConfig.SetAuthCookie(w, out.Token, out.Secret, out.XSRFToken,
			int(model.TokenExpiration().Seconds()),
			int(model.SecretExpiration.Seconds()),
			int(model.XSRFTokenExpiration.Seconds()),
			config.Config.BaseDomain)
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RegisterWithMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RefreshToken godoc
// @ID refresh-token
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.RefreshTokenResponse
// @Failure default {object} errdefs.Error
// @Router /users/refreshToken [post].
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
	if xsrfTokenHeader == "" {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token")))
		return
	}

	xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(err))
		return
	}

	secretCookie, err := r.Cookie("refresh_token")
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(err))
		return
	}

	req := requests.RefreshTokenRequest{
		Secret:          secretCookie.Value,
		XSRFTokenHeader: xsrfTokenHeader,
		XSRFTokenCookie: xsrfTokenCookie.Value,
	}
	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RefreshToken(r.Context(), adapters.RefreshTokenRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.SetAuthCookie(w, out.Token, out.Secret, out.XSRFToken,
		int(model.TokenExpiration().Seconds()),
		int(model.SecretExpiration.Seconds()),
		int(model.XSRFTokenExpiration.Seconds()),
		out.Domain)

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RefreshTokenOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SaveAuth godoc
// @ID save-auth
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SaveAuthRequest true " "
// @Success 200 {object} responses.SaveAuthResponse
// @Failure default {object} errdefs.Error
// @Router /users/saveAuth [post].
func (h *UserHandler) SaveAuth(w http.ResponseWriter, r *http.Request) {
	var req requests.SaveAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SaveAuth(r.Context(), adapters.SaveAuthRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.SetAuthCookie(w, out.Token, out.Secret, out.XSRFToken,
		int(model.TokenExpiration().Seconds()),
		int(model.SecretExpiration.Seconds()),
		int(model.XSRFTokenExpiration.Seconds()),
		out.Domain)

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.SaveAuthOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// ObtainAuthToken godoc
// @ID obtain-auth-token
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.ObtainAuthTokenResponse
// @Failure default {object} errdefs.Error
// @Router /users/obtainAuthToken [post].
func (h *UserHandler) ObtainAuthToken(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.ObtainAuthToken(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.DeleteTmpAuthCookie(w, r)

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ObtainAuthTokenOutputToResponse(out)); err != nil {
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
func (h *UserHandler) ResendInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.ResendInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.ResendInvitation(r.Context(), adapters.ResendInvitationRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.ResendInvitationOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignOut godoc
// @ID sign-out
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/signout [post].
func (h *UserHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.SignOut(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.DeleteAuthCookie(w, r, out.Domain)

	if err := httputil.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed out",
	}); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// Invite godoc
// @ID invite-users
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.InviteUsersRequest true " "
// @Success 200 {object} responses.InviteUsersResponse
// @Failure default {object} errdefs.Error
// @Router /users/invite [post].
func (h *UserHandler) Invite(w http.ResponseWriter, r *http.Request) {
	var req requests.InviteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Invite(r.Context(), adapters.InviteUsersRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.InviteUsersOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RequestInvitationMagicLink godoc
// @ID request-invitation-magic-link
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.RequestInvitationMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /users/auth/invitations/magic/request [post].
func (h *UserHandler) RequestInvitationMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestInvitationMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RequestInvitationMagicLink(r.Context(), adapters.RequestInvitationMagicLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RequestInvitationMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// AuthenticateWithInvitationMagicLink handles authentication with an invitation magic link.
func (h *UserHandler) AuthenticateWithInvitationMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.AuthenticateWithInvitationMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.AuthenticateWithInvitationMagicLink(r.Context(), adapters.AuthenticateWithInvitationMagicLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.AuthenticateWithInvitationMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RegisterWithInvitationMagicLink handles registration with an invitation magic link.
func (h *UserHandler) RegisterWithInvitationMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterWithInvitationMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RegisterWithInvitationMagicLink(r.Context(), adapters.RegisterWithInvitationMagicLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.SetAuthCookie(w, out.Token, out.Secret, out.XSRFToken,
		int(model.TokenExpiration().Seconds()),
		int(model.SecretExpiration.Seconds()),
		int(model.XSRFTokenExpiration.Seconds()),
		out.Domain)

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RegisterWithInvitationMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RequestGoogleAuthLink godoc
// @ID request-google-auth-link
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.RequestGoogleAuthLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/google/request [post].
func (h *UserHandler) RequestGoogleAuthLink(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.RequestGoogleAuthLink(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RequestGoogleAuthLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// AuthenticateWithGoogle godoc
// @ID authenticate-with-google
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.AuthenticateWithGoogleResponse
// @Failure default {object} errdefs.Error
// @Router /auth/google/authenticate [post].
func (h *UserHandler) AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req requests.AuthenticateWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrInvalidArgument(err))
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.AuthenticateWithGoogle(r.Context(), adapters.AuthenticateWithGoogleRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if !out.IsOrganizationExists && out.Flow != "invitation" {
		h.cookieConfig.SetTmpAuthCookie(w, out.Token, out.XSRFToken, config.Config.AuthDomain())
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.AuthenticateWithGoogleOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RegisterWithGoogle godoc
// @ID register-with-google
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.RegisterWithGoogleResponse
// @Failure default {object} errdefs.Error
// @Router /users/auth/google/register [post].
func (h *UserHandler) RegisterWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrInvalidArgument(err))
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RegisterWithGoogle(r.Context(), adapters.RegisterWithGoogleRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.SetTmpAuthCookie(w, out.Token, out.XSRFToken, config.Config.AuthDomain())

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RegisterWithGoogleOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RequestInvitationGoogleAuthLink godoc
// @ID request-invitation-google-auth-link
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.RequestInvitationGoogleAuthLinkResponse
// @Failure default {object} errdefs.Error
// @Router /users/auth/invitations/google/request [post].
func (h *UserHandler) RequestInvitationGoogleAuthLink(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestInvitationGoogleAuthLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrInvalidArgument(err))
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RequestInvitationGoogleAuthLink(r.Context(), adapters.RequestInvitationGoogleAuthLinkRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RequestInvitationGoogleAuthLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
