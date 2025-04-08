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
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Subject:   jwt.UserSignatureSubjectMagicLinkRegistration,
			Issuer:    jwt.Issuer,
		},
	})
}

func createGoogleAuthLinkToken(expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserGoogleAuthLinkClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expirationTime),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	})
}

func createGoogleRegistrationToken(googleID, email, firstName, lastName string, expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserGoogleRegistrationClaims{
		GoogleID:  googleID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expirationTime),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	})
}

func createGoogleInvitationStateToken(orgID uuid.UUID, expiresAt time.Time, subject string) (string, error) {
	claims := &jwt.UserGoogleInvitationRegistrationClaims{
		InvitationOrgID: orgID,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expiresAt),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	}
	return jwt.SignToken(claims)
}

func createGoogleInvitationRegistrationToken(orgID uuid.UUID, googleID, email, firstName, lastName string, expiresAt time.Time, subject string) (string, error) {
	claims := &jwt.UserGoogleInvitationRegistrationClaims{
		InvitationOrgID: orgID,
		GoogleID:        googleID,
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expiresAt),
			Issuer:    jwt.Issuer,
			Subject:   subject,
		},
	}
	return jwt.SignToken(claims)
}
