package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/auth/service"
	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
	"github.com/trysourcetool/sourcetool/backend/server/http/requests"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type AuthHandler struct {
	service      service.AuthService
	cookieConfig *CookieConfig
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service:      service,
		cookieConfig: NewCookieConfig(),
	}
}

// RequestMagicLink godoc
// @ID request-magic-link
// @Accept json
// @Produce json
// @Tags auth
// @Param Body body requests.RequestMagicLinkRequest true "Email address for magic link"
// @Success 200 {object} responses.RequestMagicLinkResponse
// @Failure 400 {object} errdefs.Error "Invalid email format"
// @Failure 404 {object} errdefs.Error "User not found"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /auth/magic/request [post].
func (h *AuthHandler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrInvalidArgument(err))
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
// @Tags auth
// @Param Body body requests.AuthenticateWithMagicLinkRequest true " "
// @Success 200 {object} responses.AuthenticateWithMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/magic/authenticate [post].
func (h *AuthHandler) AuthenticateWithMagicLink(w http.ResponseWriter, r *http.Request) {
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

	if !out.HasOrganization {
		h.cookieConfig.SetTmpAuthCookie(w, out.Token, out.XSRFToken, config.Config.AuthDomain())
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.AuthenticateWithMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RegisterWithMagicLink godoc
// @ID register-with-magic-link
// @Accept json
// @Produce json
// @Tags auth
// @Param Body body requests.RegisterWithMagicLinkRequest true "Registration data with magic link token"
// @Success 200 {object} responses.RegisterWithMagicLinkResponse
// @Failure 400 {object} errdefs.Error "Invalid request parameters"
// @Failure 401 {object} errdefs.Error "Invalid or expired magic link token"
// @Failure 500 {object} errdefs.Error "Internal server error"
// @Router /auth/magic/register [post].
func (h *AuthHandler) RegisterWithMagicLink(w http.ResponseWriter, r *http.Request) {
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
		h.cookieConfig.SetAuthCookie(w, out.Token, out.RefreshToken, out.XSRFToken,
			int(user.TokenExpiration().Seconds()),
			int(user.RefreshTokenExpiration.Seconds()),
			int(user.XSRFTokenExpiration.Seconds()),
			config.Config.BaseDomain)
	}

	if err := httputil.WriteJSON(w, http.StatusOK, adapters.RegisterWithMagicLinkOutputToResponse(out)); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

// RequestInvitationMagicLink godoc
// @ID request-invitation-magic-link
// @Accept json
// @Produce json
// @Tags auth
// @Success 200 {object} responses.RequestInvitationMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/invitations/magic/request [post].
func (h *AuthHandler) RequestInvitationMagicLink(w http.ResponseWriter, r *http.Request) {
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

// AuthenticateWithInvitationMagicLink godoc
// @ID authenticate-with-invitation-magic-link
// @Accept json
// @Produce json
// @Tags auth
// @Param Body body requests.AuthenticateWithInvitationMagicLinkRequest true " "
// @Success 200 {object} responses.AuthenticateWithInvitationMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/invitations/magic/authenticate [post].
func (h *AuthHandler) AuthenticateWithInvitationMagicLink(w http.ResponseWriter, r *http.Request) {
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

// RegisterWithInvitationMagicLink godoc
// @ID register-with-invitation-magic-link
// @Accept json
// @Produce json
// @Tags auth
// @Param Body body requests.RegisterWithInvitationMagicLinkRequest true " "
// @Success 200 {object} responses.RegisterWithInvitationMagicLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/invitations/magic/register [post].
func (h *AuthHandler) RegisterWithInvitationMagicLink(w http.ResponseWriter, r *http.Request) {
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

	h.cookieConfig.SetAuthCookie(w, out.Token, out.RefreshToken, out.XSRFToken,
		int(user.TokenExpiration().Seconds()),
		int(user.RefreshTokenExpiration.Seconds()),
		int(user.XSRFTokenExpiration.Seconds()),
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
// @Tags auth
// @Success 200 {object} responses.RequestGoogleAuthLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/google/request [post].
func (h *AuthHandler) RequestGoogleAuthLink(w http.ResponseWriter, r *http.Request) {
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
// @Tags auth
// @Success 200 {object} responses.AuthenticateWithGoogleResponse
// @Failure default {object} errdefs.Error
// @Router /auth/google/authenticate [post].
func (h *AuthHandler) AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request) {
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

	if !out.HasOrganization && out.Flow != "invitation" {
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
// @Tags auth
// @Success 200 {object} responses.RegisterWithGoogleResponse
// @Failure default {object} errdefs.Error
// @Router /auth/google/register [post].
func (h *AuthHandler) RegisterWithGoogle(w http.ResponseWriter, r *http.Request) {
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
// @Tags auth
// @Success 200 {object} responses.RequestInvitationGoogleAuthLinkResponse
// @Failure default {object} errdefs.Error
// @Router /auth/invitations/google/request [post].
func (h *AuthHandler) RequestInvitationGoogleAuthLink(w http.ResponseWriter, r *http.Request) {
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

// RefreshToken godoc
// @ID refresh-token
// @Accept json
// @Produce json
// @Tags auth
// @Success 200 {object} responses.RefreshTokenResponse
// @Failure default {object} errdefs.Error
// @Router /auth/refresh [post].
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
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

	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, errdefs.ErrUnauthenticated(err))
		return
	}

	req := requests.RefreshTokenRequest{
		RefreshToken:    refreshTokenCookie.Value,
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

	h.cookieConfig.SetAuthCookie(w, out.Token, out.RefreshToken, out.XSRFToken,
		int(user.TokenExpiration().Seconds()),
		int(user.RefreshTokenExpiration.Seconds()),
		int(user.XSRFTokenExpiration.Seconds()),
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
// @Tags auth
// @Param Body body requests.SaveAuthRequest true " "
// @Success 200 {object} responses.SaveAuthResponse
// @Failure default {object} errdefs.Error
// @Router /auth/save [post].
func (h *AuthHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req requests.SaveAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	if err := httputil.ValidateRequest(req); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	out, err := h.service.Save(r.Context(), adapters.SaveAuthRequestToDTOInput(req))
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.SetAuthCookie(w, out.Token, out.RefreshToken, out.XSRFToken,
		int(user.TokenExpiration().Seconds()),
		int(user.RefreshTokenExpiration.Seconds()),
		int(user.XSRFTokenExpiration.Seconds()),
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
// @Tags auth
// @Success 200 {object} responses.ObtainAuthTokenResponse
// @Failure default {object} errdefs.Error
// @Router /auth/token/obtain [post].
func (h *AuthHandler) ObtainAuthToken(w http.ResponseWriter, r *http.Request) {
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

// Logout godoc
// @ID logout
// @Accept json
// @Produce json
// @Tags auth
// @Success 200 {object} responses.StatusResponse
// @Failure default {object} errdefs.Error
// @Router /auth/logout [post].
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Logout(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	h.cookieConfig.DeleteAuthCookie(w, r, out.Domain)

	if err := httputil.WriteJSON(w, http.StatusOK, &responses.StatusResponse{
		Code:    http.StatusOK,
		Message: "Successfully logged out",
	}); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
