package requests

type RequestMagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthenticateWithMagicLinkRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RegisterWithMagicLinkRequest struct {
	Token     string `json:"token"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RequestInvitationMagicLinkRequest struct {
	InvitationToken string `json:"invitationToken" validate:"required"`
}

type AuthenticateWithInvitationMagicLinkRequest struct {
	Token string `json:"token" validate:"required"`
}

type RegisterWithInvitationMagicLinkRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type AuthenticateWithGoogleRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type RegisterWithGoogleRequest struct {
	Token string `json:"token" validate:"required"`
}

type RequestInvitationGoogleAuthLinkRequest struct {
	InvitationToken string `json:"invitationToken" validate:"required"`
}

type AuthenticateWithInvitationGoogleAuthLinkRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type RegisterWithInvitationGoogleAuthLinkRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken    string `validate:"required"`
	XSRFTokenHeader string `validate:"required"`
	XSRFTokenCookie string `validate:"required"`
}

type SaveAuthRequest struct {
	Token string `json:"token" validate:"required"`
}
