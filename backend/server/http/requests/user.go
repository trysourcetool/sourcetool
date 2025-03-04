package requests

type SendSignUpInstructionsRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInWithGoogleRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}

type SignUpRequest struct {
	Token                string `json:"token" validate:"required"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
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

type SignInInvitationRequest struct {
	InvitationToken string `json:"invitationToken" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type SignUpInvitationRequest struct {
	InvitationToken      string `json:"invitationToken" validate:"required"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
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

type UpdateUserPasswordRequest struct {
	CurrentPassword      string `json:"currentPassword" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
}

type SendUpdateUserEmailInstructionsRequest struct {
	Email             string `json:"email" validate:"required,email"`
	EmailConfirmation string `json:"emailConfirmation" validate:"required,email"`
}

type UpdateUserEmailRequest struct {
	Token string `json:"token" validate:"required"`
}
