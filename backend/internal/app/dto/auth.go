package dto

type RequestMagicLinkInput struct {
	Email string
}

type AuthenticateWithMagicLinkInput struct {
	Token     string
	FirstName string // Optional: used for new users
	LastName  string // Optional: used for new users
}

type RegisterWithMagicLinkInput struct {
	Token     string
	FirstName string
	LastName  string
}

type RequestInvitationMagicLinkInput struct {
	InvitationToken string
}

type AuthenticateWithInvitationMagicLinkInput struct {
	Token string
}

type RegisterWithInvitationMagicLinkInput struct {
	Token     string
	FirstName string
	LastName  string
}

type AuthenticateWithGoogleInput struct {
	Code  string
	State string
}

type RegisterWithGoogleInput struct {
	Token     string
	FirstName string
	LastName  string
}

type RequestInvitationGoogleAuthLinkInput struct {
	InvitationToken string
}

type RefreshTokenInput struct {
	RefreshToken    string
	XSRFTokenHeader string
	XSRFTokenCookie string
}

type SaveAuthInput struct {
	Token string
}

type RequestMagicLinkOutput struct {
	Email string
	IsNew bool
}

type AuthenticateWithMagicLinkOutput struct {
	AuthURL         string
	Token           string
	HasOrganization bool
	RefreshToken    string
	XSRFToken       string
	Domain          string
	IsNewUser       bool
}

type RegisterWithMagicLinkOutput struct {
	Token           string
	RefreshToken    string
	XSRFToken       string
	ExpiresAt       string
	HasOrganization bool
}

type RequestInvitationMagicLinkOutput struct {
	Email string
	IsNew bool
}

type AuthenticateWithInvitationMagicLinkOutput struct {
	AuthURL   string
	Token     string
	Domain    string
	IsNewUser bool
}

type RegisterWithInvitationMagicLinkOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	Domain       string
}

type RequestGoogleAuthLinkOutput struct {
	AuthURL string
}

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

type RegisterWithGoogleOutput struct {
	Token           string
	RefreshToken    string
	XSRFToken       string
	ExpiresAt       string
	AuthURL         string
	HasOrganization bool
}

type RequestInvitationGoogleAuthLinkOutput struct {
	AuthURL string
}

type LogoutOutput struct {
	Domain string
}

type RefreshTokenOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	Domain       string
}

type SaveAuthOutput struct {
	Token        string
	RefreshToken string
	XSRFToken    string
	ExpiresAt    string
	RedirectURL  string
	Domain       string
}

type ObtainAuthTokenOutput struct {
	AuthURL string
	Token   string
}
