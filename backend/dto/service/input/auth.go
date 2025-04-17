package input

// RequestMagicLinkInput is the input for requesting a magic link for passwordless auth.
type RequestMagicLinkInput struct {
	Email string
}

// AuthenticateWithMagicLinkInput is the input for authenticating with a magic link token.
type AuthenticateWithMagicLinkInput struct {
	Token     string
	FirstName string // Optional: used for new users
	LastName  string // Optional: used for new users
}

// RegisterWithMagicLinkInput is the input for registering with a magic link.
type RegisterWithMagicLinkInput struct {
	Token     string
	FirstName string
	LastName  string
}

// RequestInvitationMagicLinkInput represents the input for requesting a magic link for invitation.
type RequestInvitationMagicLinkInput struct {
	InvitationToken string
}

// AuthenticateWithInvitationMagicLinkInput represents the input for authenticating with an invitation magic link.
type AuthenticateWithInvitationMagicLinkInput struct {
	Token string
}

// RegisterWithInvitationMagicLinkInput represents the input for registering with an invitation magic link.
type RegisterWithInvitationMagicLinkInput struct {
	Token     string
	FirstName string
	LastName  string
}

// AuthenticateWithGoogleInput defines the input for authenticating with Google via frontend callback.
type AuthenticateWithGoogleInput struct {
	Code  string
	State string
}

// RegisterWithGoogleInput defines the input for registering a new user via Google OAuth flow.
type RegisterWithGoogleInput struct {
	Token     string
	FirstName string
	LastName  string
}

// RequestInvitationGoogleAuthLinkInput is the input for requesting a Google Auth link for an invitation.
type RequestInvitationGoogleAuthLinkInput struct {
	InvitationToken string // Token identifying the specific invitation
}

// RefreshTokenInput is the input for Refresh Token operation.
type RefreshTokenInput struct {
	RefreshToken    string
	XSRFTokenHeader string
	XSRFTokenCookie string
}

// SaveAuthInput is the input for Save Auth operation.
type SaveAuthInput struct {
	Token string
}
