package service

import (
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/jwt"
	"github.com/trysourcetool/sourcetool/backend/model"
)

func createUpdateEmailToken(userID, email string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(model.EmailTokenExpiration)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectUpdateEmail,
		},
	})
}

func createInvitationToken(email string) (string, error) {
	return jwt.SignToken(&jwt.UserClaims{
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(model.EmailTokenExpiration)),
			Issuer:    jwt.Issuer,
			Subject:   jwt.UserSignatureSubjectInvitation,
		},
	})
}
