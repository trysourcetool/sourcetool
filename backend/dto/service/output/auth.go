package output

// RequestMagicLinkOutput is the output for the magic link request operation.
type RequestMagicLinkOutput struct {
	Email string
	IsNew bool // Indicates if this is a new user
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

// RegisterWithMagicLinkOutput is the output for registering with a magic link.
type RegisterWithMagicLinkOutput struct {
	Token           string
	RefreshToken    string
	XSRFToken       string
	ExpiresAt       string
	HasOrganization bool
}

// RequestInvitationMagicLinkOutput represents the output for requesting a magic link for invitation.
type RequestInvitationMagicLinkOutput struct {
	Email string
	IsNew bool
}

// AuthenticateWithInvitationMagicLinkOutput represents the output for authenticating with an invitation magic link.
type AuthenticateWithInvitationMagicLinkOutput struct {
	AuthURL   string
	Token     string
	Domain    string
	IsNewUser bool
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

// RegisterWithGoogleOutput defines the output after successfully registering a new user via Google.
type RegisterWithGoogleOutput struct {
	Token           string
	RefreshToken    string
	XSRFToken       string
	ExpiresAt       string
	AuthURL         string
	HasOrganization bool
}

// RequestInvitationGoogleAuthLinkOutput is the output containing the Google Auth URL for an invitation.
type RequestInvitationGoogleAuthLinkOutput struct {
	AuthURL string // The URL the user should be redirected to for Google authentication
}

// LogoutOutput is the output for Logout operation.
type LogoutOutput struct {
	Domain string
}

// RefreshTokenOutput is the output for Refresh Token operation.
type RefreshTokenOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	Domain       string
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
