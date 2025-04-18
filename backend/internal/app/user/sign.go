package user

import (
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/auth"
	"github.com/trysourcetool/sourcetool/backend/jwt"
)

func createUpdateEmailToken(userID, email string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(auth.EmailTokenExpiration)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectUpdateEmail,
		},
	})
}

func createInvitationToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(auth.EmailTokenExpiration)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectInvitation,
		},
	})
}
