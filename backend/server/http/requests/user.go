package requests

type RequestMagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthenticateWithMagicLinkRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type SignInWithGoogleRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}

type SignUpWithGoogleRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
}

type RefreshTokenRequest struct {
	Secret          string `validate:"required"`
	XSRFTokenHeader string `validate:"required"`
	XSRFTokenCookie string `validate:"required"`
}

type SaveAuthRequest struct {
	Token string `json:"token" validate:"required"`
}

type InviteUsersRequest struct {
	Emails []string `json:"emails" validate:"required"`
	Role   string   `json:"role" validate:"required,oneof=admin developer member"`
}

type GoogleOAuthCallbackRequest struct {
	State string `validate:"required"`
	Code  string `validate:"required"`
}

type GetGoogleAuthCodeURLInvitationRequest struct {
	InvitationToken string `json:"invitationToken" validate:"required"`
}

type SignInWithGoogleInvitationRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}

type SignUpWithGoogleInvitationRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type SendUpdateUserEmailInstructionsRequest struct {
	Email             string `json:"email" validate:"required,email"`
	EmailConfirmation string `json:"emailConfirmation" validate:"required,email"`
}

type UpdateUserEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendInvitationRequest struct {
	InvitationID string `json:"invitationId" validate:"required,uuid"`
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

// AuthenticateWithGoogleRequest defines the request body for authenticating with Google via frontend callback.
type AuthenticateWithGoogleRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

// RegisterWithGoogleRequest defines the request body for registering a new user via Google OAuth flow.
type RegisterWithGoogleRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}
