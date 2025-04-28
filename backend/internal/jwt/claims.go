package jwt

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

const issuer = "https://auth.trysourcetool.com"

type AuthClaims struct {
	XSRFToken string
	jwt.RegisteredClaims
}

type UpdateUserEmailClaims struct {
	Email string
	jwt.RegisteredClaims
}

type InvitationClaims struct {
	jwt.RegisteredClaims
}

type MagicLinkClaims struct {
	jwt.RegisteredClaims
}

type InvitationMagicLinkClaims struct {
	jwt.RegisteredClaims
}

type MagicLinkRegistrationClaims struct {
	jwt.RegisteredClaims
}

type GoogleAuthFlow string

const (
	GoogleAuthFlowStandard   GoogleAuthFlow = "standard"
	GoogleAuthFlowInvitation GoogleAuthFlow = "invitation"
)

type GoogleAuthLinkClaims struct {
	Flow            GoogleAuthFlow
	InvitationOrgID uuid.UUID
	HostSubdomain   string
	jwt.RegisteredClaims
}

type GoogleRegistrationClaims struct {
	GoogleID        string
	FirstName       string
	LastName        string
	Flow            GoogleAuthFlow
	InvitationOrgID uuid.UUID
	Role            string // Only for invitation flow
	jwt.RegisteredClaims
}
