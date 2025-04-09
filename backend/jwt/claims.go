package jwt

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer = "trysourcetool.com"

	UserSignatureSubjectEmail                 = "email"
	UserSignatureSubjectUpdateEmail           = "update_email"
	UserSignatureSubjectActivate              = "activate"
	UserSignatureSubjectInvitation            = "invitation"
	UserSignatureSubjectInvitationMagicLink   = "invitation_magic_link"
	UserSignatureSubjectMagicLink             = "magic_link"
	UserSignatureSubjectMagicLinkRegistration = "magic_link_registration"
	UserSignatureSubjectGoogleAuthLink        = "google_auth_link"
	UserSignatureSubjectGoogleRegistration    = "google_registration"
)

type RegisteredClaims jwt.RegisteredClaims

// UserClaims represents claims for general user authentication.
type UserClaims struct {
	UserID string
	Email  string
	jwt.RegisteredClaims
}

// UserEmailClaims represents claims for email-related operations.
type UserEmailClaims struct {
	Email string
	jwt.RegisteredClaims
}

// UserAuthClaims represents claims for user authentication with XSRF token.
type UserAuthClaims struct {
	UserID    string
	XSRFToken string
	jwt.RegisteredClaims
}

// UserMagicLinkRegistrationClaims represents claims for magic link registration.
type UserMagicLinkRegistrationClaims struct {
	Email string
	jwt.RegisteredClaims
}

type GoogleAuthFlow string

const (
	GoogleAuthFlowStandard   GoogleAuthFlow = "standard"
	GoogleAuthFlowInvitation GoogleAuthFlow = "invitation"
)

// UserGoogleAuthLinkClaims represents claims for Google authentication link.
type UserGoogleAuthLinkClaims struct {
	Flow            GoogleAuthFlow
	InvitationOrgID uuid.UUID
	jwt.RegisteredClaims
}

// UserGoogleRegistrationClaims represents claims for Google registration.
type UserGoogleRegistrationClaims struct {
	GoogleID        string
	Email           string
	FirstName       string
	LastName        string
	Flow            GoogleAuthFlow
	InvitationOrgID uuid.UUID
	Role            string // Only for invitation flow
	jwt.RegisteredClaims
}
