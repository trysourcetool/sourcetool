package dto

// RequestMagicLinkInput is the input for requesting a magic link for passwordless auth.
type RequestMagicLinkInput struct {
	Email string
}

// RequestMagicLinkOutput is the output for the magic link request operation.
type RequestMagicLinkOutput struct {
	Email string
	IsNew bool // Indicates if this is a new user
}

// AuthenticateWithMagicLinkInput is the input for authenticating with a magic link token.
type AuthenticateWithMagicLinkInput struct {
	Token     string
	FirstName string // Optional: used for new users
	LastName  string // Optional: used for new users
}

// AuthenticateWithMagicLinkOutput is the output for authenticating with a magic link token.
type AuthenticateWithMagicLinkOutput struct {
	AuthURL         string
	Token           string
	HasOrganization bool
	RefreshToken    string
	XSRFToken       string
	Domain          string
	IsNewUser       bool // Indicates if a new user was created
}

// RegisterWithMagicLinkInput is the input for registering with a magic link.
type RegisterWithMagicLinkInput struct {
	Token     string
	FirstName string
	LastName  string
}

// RegisterWithMagicLinkOutput is the output for registering with a magic link.
type RegisterWithMagicLinkOutput struct {
	Token           string
	RefreshToken    string
	XSRFToken       string
	ExpiresAt       string
	HasOrganization bool
}

// RequestInvitationMagicLinkInput represents the input for requesting a magic link for invitation.
type RequestInvitationMagicLinkInput struct {
	InvitationToken string
}

// RequestInvitationMagicLinkOutput represents the output for requesting a magic link for invitation.
type RequestInvitationMagicLinkOutput struct {
	Email string
	IsNew bool
}

// AuthenticateWithInvitationMagicLinkInput represents the input for authenticating with an invitation magic link.
type AuthenticateWithInvitationMagicLinkInput struct {
	Token string
}

// AuthenticateWithInvitationMagicLinkOutput represents the output for authenticating with an invitation magic link.
type AuthenticateWithInvitationMagicLinkOutput struct {
	AuthURL   string
	Token     string
	Domain    string
	IsNewUser bool
}

// RegisterWithInvitationMagicLinkInput represents the input for registering with an invitation magic link.
type RegisterWithInvitationMagicLinkInput struct {
	Token     string
	FirstName string
	LastName  string
}

// RegisterWithInvitationMagicLinkOutput represents the output for registering with an invitation magic link.
type RegisterWithInvitationMagicLinkOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	Domain       string
}

// RequestGoogleAuthLinkOutput represents the output for requesting a Google Auth link.
type RequestGoogleAuthLinkOutput struct {
	AuthURL string
}

// AuthenticateWithGoogleInput defines the input for authenticating with Google via frontend callback.
type AuthenticateWithGoogleInput struct {
	Code  string
	State string
}

// AuthenticateWithGoogleOutput defines the output for authenticating with Google via frontend callback.
type AuthenticateWithGoogleOutput struct {
	FirstName                string
	LastName                 string
	AuthURL                  string
	Token                    string
	HasOrganization          bool
	HasMultipleOrganizations bool
	RefreshToken             string
	XSRFToken                string
	Domain                   string
	IsNewUser                bool
	Flow                     string
}

// RegisterWithGoogleInput defines the input for registering a new user via Google OAuth flow.
type RegisterWithGoogleInput struct {
	Token     string
	FirstName string
	LastName  string
}

// RegisterWithGoogleOutput defines the output after successfully registering a new user via Google.
type RegisterWithGoogleOutput struct {
	Token           string
	RefreshToken    string
	XSRFToken       string
	ExpiresAt       string
	AuthURL         string
	HasOrganization bool
}

// RequestInvitationGoogleAuthLinkInput is the input for requesting a Google Auth link for an invitation.
type RequestInvitationGoogleAuthLinkInput struct {
	InvitationToken string // Token identifying the specific invitation
}

// RequestInvitationGoogleAuthLinkOutput is the output containing the Google Auth URL for an invitation.
type RequestInvitationGoogleAuthLinkOutput struct {
	AuthURL string // The URL the user should be redirected to for Google authentication
}

// LogoutOutput is the output for Logout operation.
type LogoutOutput struct {
	Domain string
}

// RefreshTokenInput is the input for Refresh Token operation.
type RefreshTokenInput struct {
	RefreshToken    string
	XSRFTokenHeader string
	XSRFTokenCookie string
}

// RefreshTokenOutput is the output for Refresh Token operation.
type RefreshTokenOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	Domain       string
}

// SaveAuthInput is the input for Save Auth operation.
type SaveAuthInput struct {
	Token string
}

// SaveAuthOutput is the output for Save Auth operation.
type SaveAuthOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	RedirectURL  string
	Domain       string
}

// ObtainAuthTokenOutput is the output for Obtain Auth Token operation.
type ObtainAuthTokenOutput struct {
	AuthURL string
	Token   string
}
