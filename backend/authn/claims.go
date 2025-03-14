package authn

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserSignatureSubjectEmail             = "email"
	UserSignatureSubjectUpdateEmail       = "update_email"
	UserSignatureSubjectActivate          = "activate"
	UserSignatureSubjectInvitation        = "invitaiton"
	UserSignatureSubjectGoogleAuthRequest = "google_auth_request"
)

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
	Email     string
	XSRFToken string
	jwt.RegisteredClaims
}
