package user

import (
	"time"

	"github.com/gofrs/uuid/v5"
	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/jwt"
)

func createAuthToken(userID, xsrfToken string, expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserAuthClaims{
		UserID:    userID,
		XSRFToken: xsrfToken,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expirationTime),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	})
}

func createUserToken(userID, email string, expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expirationTime),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	})
}

func createUserEmailToken(email string, expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserEmailClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expirationTime),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	})
}

func createRegistrationToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserMagicLinkRegistrationClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Subject:   jwt.UserSignatureSubjectMagicLinkRegistration,
			Issuer:    jwt.Issuer,
		},
	})
}

func createGoogleAuthLinkToken(flow jwt.GoogleAuthFlow, invitationOrgID uuid.UUID, hostSubdomain string) (string, error) {
	claims := &jwt.UserGoogleAuthLinkClaims{
		Flow:            flow,
		InvitationOrgID: invitationOrgID,
		HostSubdomain:   hostSubdomain,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectGoogleAuthLink,
		},
	}
	return jwt.SignToken(claims)
}

func createGoogleRegistrationToken(googleID, email, firstName, lastName string, flow jwt.GoogleAuthFlow, invitationOrgID uuid.UUID, role string) (string, error) {
	claims := &jwt.UserGoogleRegistrationClaims{
		GoogleID:        googleID,
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		Flow:            flow,
		InvitationOrgID: invitationOrgID,
		Role:            role,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectGoogleRegistration,
		},
	}
	return jwt.SignToken(claims)
}
