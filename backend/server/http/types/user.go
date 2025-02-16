package types

type SendSignUpInstructionsInput struct {
	Email string `json:"email" validate:"required,email"`
}

type SendSignUpInstructionsPayload struct {
	Email string `json:"email"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInPayload struct {
	AuthURL              string `json:"authUrl"`
	Token                string `json:"token"`
	IsOrganizationExists bool   `json:"isOrganizationExists"`
	Secret               string `json:"-"`
	XSRFToken            string `json:"-"`
	Domain               string `json:"-"`
}

type SignInWithGoogleInput struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}

type SignInWithGooglePayload struct {
	AuthURL              string `json:"authUrl"`
	Token                string `json:"token"`
	IsOrganizationExists bool   `json:"isOrganizationExists"`
	Secret               string `json:"-"`
	XSRFToken            string `json:"-"`
	Domain               string `json:"-"`
}

type SignUpInput struct {
	Token                string `json:"token" validate:"required"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
}

type SignUpPayload struct {
	Token     string
	XSRFToken string
}

type SignUpWithGoogleInput struct {
	SessionToken string `json:"sessionToken" validate:"required"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
}

type SignUpWithGooglePayload struct {
	Token     string
	XSRFToken string
}

type RefreshTokenInput struct {
	Secret          string `validate:"required"`
	XSRFTokenHeader string `validate:"required"`
	XSRFTokenCookie string `validate:"required"`
}

type RefreshTokenPayload struct {
	Token     string `json:"-"`
	Secret    string `json:"-"`
	XSRFToken string `json:"-"`
	ExpiresAt string `json:"expiresAt"`
	Domain    string `json:"-"`
}

type SaveAuthInput struct {
	Token string `json:"token" validate:"required"`
}

type SaveAuthPayload struct {
	Token       string `json:"-"`
	Secret      string `json:"-"`
	XSRFToken   string `json:"-"`
	ExpiresAt   string `json:"expiresAt"`
	RedirectURL string `json:"redirectUrl"`
	Domain      string `json:"-"`
}

type ObtainAuthTokenPayload struct {
	AuthURL   string `json:"authUrl"`
	Token     string `json:"token"`
	Secret    string `json:"-"`
	XSRFToken string `json:"-"`
	Domain    string `json:"-"`
}

type InviteUsersInput struct {
	Emails []string `json:"emails" validate:"required"`
	Role   string   `json:"role" validate:"required,oneof=admin developer member"`
}

type InviteUsersPayload struct {
	UserInvitations []*UserInvitationPayload `json:"userInvitations"`
}

type SignInInvitationInput struct {
	InvitationToken string `json:"invitationToken" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type SignInInvitationPayload struct {
	Token     string `json:"-"`
	Secret    string `json:"-"`
	XSRFToken string `json:"-"`
	ExpiresAt string `json:"expiresAt"`
	Domain    string `json:"-"`
}

type SignUpInvitationInput struct {
	InvitationToken      string `json:"invitationToken" validate:"required"`
	FirstName            string `json:"firstName" validate:"required"`
	LastName             string `json:"lastName" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
}

type SignUpInvitationPayload struct {
	Token     string `json:"-"`
	Secret    string `json:"-"`
	XSRFToken string `json:"-"`
	ExpiresAt string `json:"expiresAt"`
	Domain    string `json:"-"`
}

type GetGoogleAuthCodeURLPayload struct {
	URL string `json:"url"`
}

type GoogleOAuthCallbackInput struct {
	State string `validate:"required"`
	Code  string `validate:"required"`
}

type GoogleOAuthCallbackPayload struct {
	SessionToken string
	IsUserExists bool
	FirstName    string
	LastName     string
	Domain       string
	Invited      bool
}

type GetGoogleAuthCodeURLInvitationInput struct {
	InvitationToken string `json:"invitationToken" validate:"required"`
}

type GetGoogleAuthCodeURLInvitationPayload struct {
	URL string `json:"url"`
}

type SignInWithGoogleInvitationInput struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}

type SignInWithGoogleInvitationPayload struct {
	Token     string `json:"-"`
	Secret    string `json:"-"`
	XSRFToken string `json:"-"`
	ExpiresAt string `json:"expiresAt"`
	Domain    string `json:"-"`
}

type SignUpWithGoogleInvitationInput struct {
	SessionToken string `json:"sessionToken" validate:"required"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
}

type SignUpWithGoogleInvitationPayload struct {
	Token     string `json:"-"`
	Secret    string `json:"-"`
	XSRFToken string `json:"-"`
	ExpiresAt string `json:"expiresAt"`
	Domain    string `json:"-"`
}

type SignOutPayload struct {
	Domain string
}

type UserPayload struct {
	ID           string               `json:"id"`
	Email        string               `json:"email"`
	FirstName    string               `json:"firstName"`
	LastName     string               `json:"lastName"`
	Role         string               `json:"role"`
	CreatedAt    string               `json:"createdAt"`
	UpdatedAt    string               `json:"updatedAt"`
	Organization *OrganizationPayload `json:"organization"`
}

type UserInvitationPayload struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type ListUsersPayload struct {
	Users           []*UserPayload           `json:"users"`
	UserInvitations []*UserInvitationPayload `json:"userInvitations"`
}

type GetMePayload struct {
	User *UserPayload `json:"user"`
}

type UpdateUserInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type UpdateUserPayload struct {
	User *UserPayload `json:"user"`
}

type UpdateUserPasswordInput struct {
	CurrentPassword      string `json:"currentPassword" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required"`
}

type UpdateUserPasswordPayload struct {
	User *UserPayload `json:"user"`
}

type SendUpdateUserEmailInstructionsInput struct {
	Email             string `json:"email" validate:"required,email"`
	EmailConfirmation string `json:"emailConfirmation" validate:"required,email"`
}

type UpdateUserEmailInput struct {
	Token string `json:"token" validate:"required"`
}

type UpdateUserEmailPayload struct {
	User *UserPayload `json:"user"`
}

type UserGroupPayload struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	GroupID   string `json:"groupId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
