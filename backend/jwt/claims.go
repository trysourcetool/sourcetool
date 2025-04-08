package jwt

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer = "trysourcetool.com"

	UserSignatureSubjectEmail                        = "email"
	UserSignatureSubjectUpdateEmail                  = "update_email"
	UserSignatureSubjectActivate                     = "activate"
	UserSignatureSubjectInvitation                   = "invitation"
	UserSignatureSubjectInvitationMagicLink          = "invitation_magic_link"
	UserSignatureSubjectGoogleAuthRequest            = "google_auth_request"
	UserSignatureSubjectGoogleRegistration           = "google_registration"
	UserSignatureSubjectMagicLink                    = "magic_link"
	UserSignatureSubjectMagicLinkRegistration        = "magic_link_registration"
	UserSignatureSubjectGoogleAuthLink               = "google_auth_link"
	UserSignatureSubjectGoogleAuthLinkInvitation     = "google_auth_link_invitation"
	UserSignatureSubjectGoogleInvitationRegistration = "google_invitation_registration"
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

// UserGoogleAuthLinkClaims represents claims for Google authentication link.
type UserGoogleAuthLinkClaims struct {
	jwt.RegisteredClaims
}

// UserGoogleRegistrationClaims represents claims for Google registration.
type UserGoogleRegistrationClaims struct {
	GoogleID  string
	Email     string
	FirstName string
	LastName  string
	jwt.RegisteredClaims
}

// UserGoogleInvitationRegistrationClaims represents claims for Google invitation registration.
type UserGoogleInvitationRegistrationClaims struct {
	InvitationOrgID uuid.UUID
	GoogleID        string
	Email           string
	FirstName       string
	LastName        string
	jwt.RegisteredClaims
}
