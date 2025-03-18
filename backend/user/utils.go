package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"path"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/utils/urlutil"
)

func hashSecret(plainSecret string) string {
	hash := sha256.Sum256([]byte(plainSecret))
	return hex.EncodeToString(hash[:])
}

func generateSecret() (plainSecret, hashedSecret string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainSecret = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	hashedSecret = hashSecret(plainSecret)

	return plainSecret, hashedSecret, nil
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
