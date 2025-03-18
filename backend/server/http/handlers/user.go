package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/httputils"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service}
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.GetMeOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.ListUsersOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Update(r.Context(), adapters.UpdateUserRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.UpdateUserOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.SendUpdateEmailInstructions(r.Context(), adapters.SendUpdateUserEmailInstructionsRequestToDTOInput(req)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully sent update email instructions",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateEmail(r.Context(), adapters.UpdateUserEmailRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.UpdateUserEmailOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdatePassword godoc
// @ID update-user-password
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.UpdateUserPasswordRequest true " "
// @Success 200 {object} responses.UpdateUserPasswordResponse
// @Failure default {object} errdefs.Error
// @Router /users/password [put].
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateUserPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdatePassword(r.Context(), adapters.UpdateUserPasswordRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.UpdateUserPasswordOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignIn godoc
// @ID signin
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignInRequest true " "
// @Success 200 {object} responses.SignInResponse
// @Failure default {object} errdefs.Error
// @Router /users/signin [post].
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req requests.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignIn(r.Context(), adapters.SignInRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if !out.IsOrganizationExists {
		maxAge := int(model.TmpTokenExpiration.Seconds())
		h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, maxAge, maxAge, maxAge)
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.SignInOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignInWithGoogle godoc
// @ID signin-with-google
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignInWithGoogleRequest true " "
// @Success 200 {object} responses.SignInWithGoogleResponse
// @Failure default {object} errdefs.Error
// @Router /users/oauth/google/signin [post].
func (h *UserHandler) SignInWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req requests.SignInWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignInWithGoogle(r.Context(), adapters.SignInWithGoogleRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if !out.IsOrganizationExists {
		maxAge := int(model.TmpTokenExpiration.Seconds())
		h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, maxAge, maxAge, maxAge)
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.SignInWithGoogleOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SendSignUpInstructions godoc
// @ID signup-instructions
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SendSignUpInstructionsRequest true " "
// @Success 200 {object} responses.SendSignUpInstructionsResponse
// @Failure default {object} errdefs.Error
// @Router /users/signup/instructions [post].
func (h *UserHandler) SendSignUpInstructions(w http.ResponseWriter, r *http.Request) {
	var req requests.SendSignUpInstructionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SendSignUpInstructions(r.Context(), adapters.SendSignUpInstructionsRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.SendSignUpInstructionsOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignUp godoc
// @ID signup
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignUpRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/signup [post].
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUp(r.Context(), adapters.SignUpRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setTmpAuthCookie(w, out.Token, out.XSRFToken)

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed up",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignUpWithGoogle godoc
// @ID signup-with-google
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignUpWithGoogleRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/oauth/google/signup [post].
func (h *UserHandler) SignUpWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUpWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUpWithGoogle(r.Context(), adapters.SignUpWithGoogleRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setTmpAuthCookie(w, out.Token, out.XSRFToken)

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed up with Google",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token")))
		return
	}

	xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(err))
		return
	}

	secretCookie, err := r.Cookie("refresh_token")
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(err))
		return
	}

	req := requests.RefreshTokenRequest{
		Secret:          secretCookie.Value,
		XSRFTokenHeader: xsrfTokenHeader,
		XSRFTokenCookie: xsrfTokenCookie.Value,
	}
	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RefreshToken(r.Context(), adapters.RefreshTokenRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.RefreshTokenOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SaveAuth(r.Context(), adapters.SaveAuthRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.SaveAuthOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.deleteTmpAuthCookie(w, r)

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.ResendInvitation(r.Context(), adapters.ResendInvitationRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.ResendInvitationOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.deleteAuthCookie(w, r, out.Domain)

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed out",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
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
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Invite(r.Context(), adapters.InviteUsersRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.InviteUsersOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignInInvitation godoc
// @ID signin-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignInInvitationRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/signin [post].
func (h *UserHandler) SignInInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.SignInInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignInInvitation(r.Context(), adapters.SignInInvitationRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed in",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignUpInvitation godoc
// @ID signup-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignUpInvitationRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/signup [post].
func (h *UserHandler) SignUpInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUpInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUpInvitation(r.Context(), adapters.SignUpInvitationRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed up",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// GetGoogleAuthCodeURL godoc
// @ID get-google-auth-code-url
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} responses.GetGoogleAuthCodeURLResponse
// @Failure default {object} errdefs.Error
// @Router /users/oauth/google/authCodeUrl [post].
func (h *UserHandler) GetGoogleAuthCodeURL(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.GetGoogleAuthCodeURL(r.Context())
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

func (h *UserHandler) GoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	req := requests.GoogleOAuthCallbackRequest{
		State: r.URL.Query().Get("state"),
		Code:  r.URL.Query().Get("code"),
	}
	if err := httputils.ValidateRequest(req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	out, err := h.service.GoogleOAuthCallback(r.Context(), adapters.GoogleOAuthCallbackRequestToDTOInput(req))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	path := "/users/oauth/google/callback"
	if out.Invited {
		path = "/users/oauth/google/invitations/callback"
	}
	base := httputils.HTTPScheme() + "://" + out.Domain + path
	params := url.Values{}
	params.Add("token", out.SessionToken)
	params.Add("isUserExists", strconv.FormatBool(out.IsUserExists))
	if out.FirstName != "" {
		params.Add("firstName", out.FirstName)
	}
	if out.LastName != "" {
		params.Add("lastName", out.LastName)
	}
	targetURL := base + "?" + params.Encode()
	http.Redirect(w, r, targetURL, http.StatusFound)
}

// GetGoogleAuthCodeURLInvitation godoc
// @ID get-google-auth-code-url-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.GetGoogleAuthCodeURLInvitationRequest true " "
// @Success 200 {object} responses.GetGoogleAuthCodeURLInvitationResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/oauth/google/authCodeUrl [post].
func (h *UserHandler) GetGoogleAuthCodeURLInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.GetGoogleAuthCodeURLInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.GetGoogleAuthCodeURLInvitation(r.Context(), adapters.GetGoogleAuthCodeURLInvitationRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, adapters.GetGoogleAuthCodeURLInvitationOutputToResponse(out)); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignInWithGoogleInvitation godoc
// @ID signin-with-google-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignInWithGoogleInvitationRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/oauth/google/signin [post].
func (h *UserHandler) SignInWithGoogleInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.SignInWithGoogleInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignInWithGoogleInvitation(r.Context(), adapters.SignInWithGoogleInvitationRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed in with Google",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignUpWithGoogleInvitation godoc
// @ID signup-with-google-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body requests.SignUpWithGoogleInvitationRequest true " "
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /users/invitations/oauth/google/signup [post].
func (h *UserHandler) SignUpWithGoogleInvitation(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUpWithGoogleInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(req); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUpWithGoogleInvitation(r.Context(), adapters.SignUpWithGoogleInvitationRequestToDTOInput(req))
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully signed up with Google",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

func (h *UserHandler) setTmpAuthCookie(w http.ResponseWriter, token, xsrfToken string) {
	xsrfTokenSameSite := http.SameSiteNoneMode
	if config.Config.Env == config.EnvLocal {
		xsrfTokenSameSite = http.SameSiteLaxMode
	}
	domain := config.Config.AuthDomain()
	maxAge := int(model.TmpTokenExpiration.Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "xsrf_token",
		Value:    xsrfToken,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: false,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: xsrfTokenSameSite,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "xsrf_token_same_site",
		Value:    xsrfToken,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *UserHandler) deleteTmpAuthCookie(w http.ResponseWriter, r *http.Request) {
	xsrfTokenSameSite := http.SameSiteNoneMode
	if config.Config.Env == config.EnvLocal {
		xsrfTokenSameSite = http.SameSiteLaxMode
	}
	domain := config.Config.AuthDomain()
	tokenCookie, _ := r.Cookie("access_token")
	if tokenCookie != nil {
		tokenCookie.MaxAge = -1
		tokenCookie.Domain = domain
		tokenCookie.Path = "/"
		tokenCookie.HttpOnly = true
		tokenCookie.Secure = !(config.Config.Env == config.EnvLocal)
		tokenCookie.SameSite = http.SameSiteStrictMode
		http.SetCookie(w, tokenCookie)
	}
	xsrfTokenCookie, _ := r.Cookie("xsrf_token")
	if xsrfTokenCookie != nil {
		xsrfTokenCookie.MaxAge = -1
		xsrfTokenCookie.Domain = domain
		xsrfTokenCookie.Path = "/"
		xsrfTokenCookie.HttpOnly = false
		xsrfTokenCookie.Secure = !(config.Config.Env == config.EnvLocal)
		xsrfTokenCookie.SameSite = xsrfTokenSameSite
		http.SetCookie(w, xsrfTokenCookie)
	}
	xsrfTokenSameSiteCookie, _ := r.Cookie("xsrf_token_same_site")
	if xsrfTokenSameSiteCookie != nil {
		xsrfTokenSameSiteCookie.MaxAge = -1
		xsrfTokenSameSiteCookie.Domain = domain
		xsrfTokenSameSiteCookie.Path = "/"
		xsrfTokenSameSiteCookie.HttpOnly = true
		xsrfTokenSameSiteCookie.Secure = !(config.Config.Env == config.EnvLocal)
		xsrfTokenSameSiteCookie.SameSite = http.SameSiteStrictMode
		http.SetCookie(w, xsrfTokenSameSiteCookie)
	}
}

func (h *UserHandler) setAuthCookie(w http.ResponseWriter, domain, token, secret, xsrfToken string, tokenMaxAge, secretMaxAge, xsrfTokenMaxAge int) {
	xsrfTokenSameSite := http.SameSiteNoneMode
	if config.Config.Env == config.EnvLocal {
		xsrfTokenSameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   tokenMaxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    secret,
		MaxAge:   secretMaxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "xsrf_token",
		Value:    xsrfToken,
		MaxAge:   xsrfTokenMaxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: false,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: xsrfTokenSameSite,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "xsrf_token_same_site",
		Value:    xsrfToken,
		MaxAge:   xsrfTokenMaxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   !(config.Config.Env == config.EnvLocal),
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *UserHandler) deleteAuthCookie(w http.ResponseWriter, r *http.Request, domain string) {
	xsrfTokenSameSite := http.SameSiteNoneMode
	if config.Config.Env == config.EnvLocal {
		xsrfTokenSameSite = http.SameSiteLaxMode
	}
	tokenCookie, _ := r.Cookie("access_token")
	if tokenCookie != nil {
		tokenCookie.MaxAge = -1
		tokenCookie.Domain = domain
		tokenCookie.Path = "/"
		tokenCookie.HttpOnly = true
		tokenCookie.Secure = !(config.Config.Env == config.EnvLocal)
		tokenCookie.SameSite = http.SameSiteStrictMode
		http.SetCookie(w, tokenCookie)
	}
	secretCookie, _ := r.Cookie("refresh_token")
	if secretCookie != nil {
		secretCookie.MaxAge = -1
		secretCookie.Domain = domain
		secretCookie.Path = "/"
		secretCookie.HttpOnly = true
		secretCookie.Secure = !(config.Config.Env == config.EnvLocal)
		secretCookie.SameSite = http.SameSiteStrictMode
		http.SetCookie(w, secretCookie)
	}
	xsrfTokenCookie, _ := r.Cookie("xsrf_token")
	if xsrfTokenCookie != nil {
		xsrfTokenCookie.MaxAge = -1
		xsrfTokenCookie.Domain = domain
		xsrfTokenCookie.Path = "/"
		xsrfTokenCookie.HttpOnly = false
		xsrfTokenCookie.Secure = !(config.Config.Env == config.EnvLocal)
		xsrfTokenCookie.SameSite = xsrfTokenSameSite
		http.SetCookie(w, xsrfTokenCookie)
	}
	xsrfTokenSameSiteCookie, _ := r.Cookie("xsrf_token_same_site")
	if xsrfTokenSameSiteCookie != nil {
		xsrfTokenSameSiteCookie.MaxAge = -1
		xsrfTokenSameSiteCookie.Domain = domain
		xsrfTokenSameSiteCookie.Path = "/"
		xsrfTokenSameSiteCookie.HttpOnly = true
		xsrfTokenSameSiteCookie.Secure = !(config.Config.Env == config.EnvLocal)
		xsrfTokenSameSiteCookie.SameSite = http.SameSiteStrictMode
		http.SetCookie(w, xsrfTokenSameSiteCookie)
	}
}
