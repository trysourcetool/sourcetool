package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer = "trysourcetool.com"

	UserSignatureSubjectEmail                 = "email"
	UserSignatureSubjectUpdateEmail           = "update_email"
	UserSignatureSubjectActivate              = "activate"
	UserSignatureSubjectInvitation            = "invitaiton"
	UserSignatureSubjectGoogleAuthRequest     = "google_auth_request"
	UserSignatureSubjectMagicLink             = "magic_link"
	UserSignatureSubjectMagicLinkRegistration = "magic_link_registration"
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

// UserGoogleAuthRequestClaims represents claims for Google authentication.
type UserGoogleAuthRequestClaims struct {
	GoogleAuthRequestID string
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
