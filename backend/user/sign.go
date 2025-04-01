package user

import (
	"time"

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

func createGoogleAuthRequestToken(googleAuthRequestID string, expirationTime time.Time, subject string) (string, error) {
	return jwt.SignToken(&jwt.UserGoogleAuthRequestClaims{
		GoogleAuthRequestID: googleAuthRequestID,
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
