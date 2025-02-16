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
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
	"github.com/trysourcetool/sourcetool/backend/user"
)

type UserHandlerCE interface {
	GetMe(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	SendUpdateEmailInstructions(w http.ResponseWriter, r *http.Request)
	UpdateEmail(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignInWithGoogle(w http.ResponseWriter, r *http.Request)
	SendSignUpInstructions(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	SignUpWithGoogle(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	SaveAuth(w http.ResponseWriter, r *http.Request)
	ObtainAuthToken(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	Invite(w http.ResponseWriter, r *http.Request)
	SignInInvitation(w http.ResponseWriter, r *http.Request)
	SignUpInvitation(w http.ResponseWriter, r *http.Request)
	GetGoogleAuthCodeURL(w http.ResponseWriter, r *http.Request)
	GoogleOAuthCallback(w http.ResponseWriter, r *http.Request)
	GetGoogleAuthCodeURLInvitation(w http.ResponseWriter, r *http.Request)
	SignInWithGoogleInvitation(w http.ResponseWriter, r *http.Request)
	SignUpWithGoogleInvitation(w http.ResponseWriter, r *http.Request)
}

type UserHandlerCEImpl struct {
	service user.ServiceCE
}

func NewUserHandlerCE(service user.ServiceCE) *UserHandlerCEImpl {
	return &UserHandlerCEImpl{service}
}

// GetMe godoc
// @ID get-me
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} types.GetMePayload
// @Failure default {object} errdefs.Error
// @Router /users/me [get].
func (h *UserHandlerCEImpl) GetMe(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.GetMe(r.Context())
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
// @ID list-users
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} types.ListUsersPayload
// @Failure default {object} errdefs.Error
// @Router /users [get].
func (h *UserHandlerCEImpl) List(w http.ResponseWriter, r *http.Request) {
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

// Update godoc
// @ID update-user
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.UpdateUserInput true " "
// @Success 200 {object} types.UpdateUserPayload
// @Failure default {object} errdefs.Error
// @Router /users [put].
func (h *UserHandlerCEImpl) Update(w http.ResponseWriter, r *http.Request) {
	var in types.UpdateUserInput
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

// SendUpdateEmailInstructions godoc
// @ID send-update-email-instructions
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SendUpdateUserEmailInstructionsInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/sendUpdateEmailInstructions [post].
func (h *UserHandlerCEImpl) SendUpdateEmailInstructions(w http.ResponseWriter, r *http.Request) {
	var in types.SendUpdateUserEmailInstructionsInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := h.service.SendUpdateEmailInstructions(r.Context(), in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Param Body body types.UpdateUserEmailInput true " "
// @Success 200 {object} types.UpdateUserEmailPayload
// @Failure default {object} errdefs.Error
// @Router /users/email [put].
func (h *UserHandlerCEImpl) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	var in types.UpdateUserEmailInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdateEmail(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// UpdatePassword godoc
// @ID update-user-password
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.UpdateUserPasswordInput true " "
// @Success 200 {object} types.UpdateUserPasswordPayload
// @Failure default {object} errdefs.Error
// @Router /users/password [put].
func (h *UserHandlerCEImpl) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var in types.UpdateUserPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.UpdatePassword(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignIn godoc
// @ID signin
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/signin [post].
func (h *UserHandlerCEImpl) SignIn(w http.ResponseWriter, r *http.Request) {
	var in types.SignInInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignIn(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if !out.IsOrganizationExists {
		maxAge := int(model.TmpTokenExpiration.Seconds())
		h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, maxAge, maxAge, maxAge)
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignInWithGoogle godoc
// @ID signin-with-google
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SignInWithGoogleInput true " "
// @Success 200 {object} types.SignInWithGooglePayload
// @Failure default {object} errdefs.Error
// @Router /users/oauth/google/signin [post].
func (h *UserHandlerCEImpl) SignInWithGoogle(w http.ResponseWriter, r *http.Request) {
	var in types.SignInWithGoogleInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignInWithGoogle(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if !out.IsOrganizationExists {
		maxAge := int(model.TmpTokenExpiration.Seconds())
		h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, maxAge, maxAge, maxAge)
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SendSignUpInstructions godoc
// @ID signup-instructions
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SendSignUpInstructionsInput true " "
// @Success 200 {object} types.SendSignUpInstructionsPayload
// @Failure default {object} errdefs.Error
// @Router /users/signup/instructions [post].
func (h *UserHandlerCEImpl) SendSignUpInstructions(w http.ResponseWriter, r *http.Request) {
	var in types.SendSignUpInstructionsInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SendSignUpInstructions(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignUp godoc
// @ID signup
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SignUpInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/signup [post].
func (h *UserHandlerCEImpl) SignUp(w http.ResponseWriter, r *http.Request) {
	var in types.SignUpInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUp(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setTmpAuthCookie(w, out.Token, out.XSRFToken)

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Param Body body types.SignUpWithGoogleInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/oauth/google/signup [post].
func (h *UserHandlerCEImpl) SignUpWithGoogle(w http.ResponseWriter, r *http.Request) {
	var in types.SignUpWithGoogleInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUpWithGoogle(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setTmpAuthCookie(w, out.Token, out.XSRFToken)

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Success 200 {object} types.RefreshTokenPayload
// @Failure default {object} errdefs.Error
// @Router /users/refreshToken [post].
func (h *UserHandlerCEImpl) RefreshToken(w http.ResponseWriter, r *http.Request) {
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

	in := types.RefreshTokenInput{
		Secret:          secretCookie.Value,
		XSRFTokenHeader: xsrfTokenHeader,
		XSRFTokenCookie: xsrfTokenCookie.Value,
	}
	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.RefreshToken(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SaveAuth godoc
// @ID save-auth
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SaveAuthInput true " "
// @Success 200 {object} types.SaveAuthPayload
// @Failure default {object} errdefs.Error
// @Router /users/saveAuth [post].
func (h *UserHandlerCEImpl) SaveAuth(w http.ResponseWriter, r *http.Request) {
	var in types.SaveAuthInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SaveAuth(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// ObtainAuthToken godoc
// @ID obtain-auth-token
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} types.ObtainAuthTokenPayload
// @Failure default {object} errdefs.Error
// @Router /users/obtainAuthToken [post].
func (h *UserHandlerCEImpl) ObtainAuthToken(w http.ResponseWriter, r *http.Request) {
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

// SignOut godoc
// @ID sign-out
// @Accept json
// @Produce json
// @Tags users
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/signout [post].
func (h *UserHandlerCEImpl) SignOut(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.SignOut(r.Context())
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.deleteAuthCookie(w, r, out.Domain)

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Param Body body types.InviteUsersInput true " "
// @Success 200 {object} types.InviteUsersPayload
// @Failure default {object} errdefs.Error
// @Router /users/invite [post].
func (h *UserHandlerCEImpl) Invite(w http.ResponseWriter, r *http.Request) {
	var in types.InviteUsersInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Invite(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignInInvitation godoc
// @ID signin-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SignInInvitationInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/invitations/signin [post].
func (h *UserHandlerCEImpl) SignInInvitation(w http.ResponseWriter, r *http.Request) {
	var in types.SignInInvitationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignInInvitation(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Param Body body types.SignUpInvitationInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/invitations/signup [post].
func (h *UserHandlerCEImpl) SignUpInvitation(w http.ResponseWriter, r *http.Request) {
	var in types.SignUpInvitationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUpInvitation(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Success 200 {object} types.GetGoogleAuthCodeURLPayload
// @Failure default {object} errdefs.Error
// @Router /users/oauth/google/authCodeUrl [post].
func (h *UserHandlerCEImpl) GetGoogleAuthCodeURL(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandlerCEImpl) GoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	in := types.GoogleOAuthCallbackInput{
		State: r.URL.Query().Get("state"),
		Code:  r.URL.Query().Get("code"),
	}
	if err := httputils.ValidateRequest(in); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	out, err := h.service.GoogleOAuthCallback(r.Context(), in)
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
// @Param Body body types.GetGoogleAuthCodeURLInvitationInput true " "
// @Success 200 {object} types.GetGoogleAuthCodeURLInvitationPayload
// @Failure default {object} errdefs.Error
// @Router /users/invitations/oauth/google/authCodeUrl [post].
func (h *UserHandlerCEImpl) GetGoogleAuthCodeURLInvitation(w http.ResponseWriter, r *http.Request) {
	var in types.GetGoogleAuthCodeURLInvitationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.GetGoogleAuthCodeURLInvitation(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, out); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// SignInWithGoogleInvitation godoc
// @ID signin-with-google-invitation
// @Accept json
// @Produce json
// @Tags users
// @Param Body body types.SignInWithGoogleInvitationInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/invitations/oauth/google/signin [post].
func (h *UserHandlerCEImpl) SignInWithGoogleInvitation(w http.ResponseWriter, r *http.Request) {
	var in types.SignInWithGoogleInvitationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignInWithGoogleInvitation(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
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
// @Param Body body types.SignUpWithGoogleInvitationInput true " "
// @Success 200 {object} types.SuccessPayload
// @Failure default {object} errdefs.Error
// @Router /users/invitations/oauth/google/signup [post].
func (h *UserHandlerCEImpl) SignUpWithGoogleInvitation(w http.ResponseWriter, r *http.Request) {
	var in types.SignUpWithGoogleInvitationInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputils.ValidateRequest(in); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.SignUpWithGoogleInvitation(r.Context(), in)
	if err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.setAuthCookie(w, out.Domain, out.Token, out.Secret, out.XSRFToken, int(model.TokenExpiration().Seconds()), int(model.SecretExpiration.Seconds()), int(model.XSRFTokenExpiration.Seconds()))

	if err := httputils.WriteJSON(w, http.StatusOK, &types.SuccessPayload{
		Code:    http.StatusOK,
		Message: "Successfully signed up with Google",
	}); err != nil {
		httputils.WriteErrJSON(r.Context(), w, err)
		return
	}
}

func (h *UserHandlerCEImpl) setTmpAuthCookie(w http.ResponseWriter, token, xsrfToken string) {
	xsrfTokenSameSite := http.SameSiteNoneMode
	if config.Config.Env == config.EnvLocal {
		xsrfTokenSameSite = http.SameSiteLaxMode
	}
	domain := "auth" + "." + config.Config.Domain
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

func (h *UserHandlerCEImpl) deleteTmpAuthCookie(w http.ResponseWriter, r *http.Request) {
	xsrfTokenSameSite := http.SameSiteNoneMode
	if config.Config.Env == config.EnvLocal {
		xsrfTokenSameSite = http.SameSiteLaxMode
	}
	domain := "auth" + "." + config.Config.Domain
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

func (h *UserHandlerCEImpl) setAuthCookie(w http.ResponseWriter, domain, token, secret, xsrfToken string, tokenMaxAge, secretMaxAge, xsrfTokenMaxAge int) {
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

func (h *UserHandlerCEImpl) deleteAuthCookie(w http.ResponseWriter, r *http.Request, domain string) {
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
