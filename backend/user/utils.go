package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"path"
	"strconv"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/jwt"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/utils/urlutil"
)

// hashSecret creates a SHA-256 hash of a plaintext secret.
func hashSecret(plainSecret string) string {
	hash := sha256.Sum256([]byte(plainSecret))
	return hex.EncodeToString(hash[:])
}

// generateSecret creates a new random secret and its hash for user authentication.
func generateSecret() (plainSecret, hashedSecret string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainSecret = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	hashedSecret = hashSecret(plainSecret)

	return plainSecret, hashedSecret, nil
}

// hashPassword creates a bcrypt hash of the password.
func hashPassword(password string) (string, error) {
	encodedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return hex.EncodeToString(encodedPass[:]), nil
}

// verifyPassword checks if the provided password matches the stored hash.
func verifyPassword(storedHash, password string) error {
	h, err := hex.DecodeString(storedHash)
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	if err = bcrypt.CompareHashAndPassword(h, []byte(password)); err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	return nil
}

func buildUserActivateURL(token string) (string, error) {
	return urlutil.BuildURL(config.Config.AuthBaseURL(), path.Join("signup", "activate"), map[string]string{
		"token": token,
	})
}

func buildUpdateEmailURL(subdomain, token string) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "email", "update", "confirm"), map[string]string{
		"token": token,
	})
}

func buildInvitationURL(subdomain, token, email string, isUserExists bool) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "invitation", "activate"), map[string]string{
		"token":        token,
		"email":        email,
		"isUserExists": strconv.FormatBool(isUserExists),
	})
}

func buildSaveAuthURL(subdomain string) (string, error) {
	return config.Config.OrgBaseURL(subdomain) + model.SaveAuthPath, nil
}

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
