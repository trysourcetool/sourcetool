package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/model"
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

func buildUserActivateURL(ctx context.Context, token string) (string, error) {
	u, err := url.Parse(ctxutils.HTTPReferer(ctx))
	if err != nil {
		return "", nil
	}

	u.Path = path.Join(u.Path, "signup", "activate")

	q := u.Query()
	q.Add("token", token)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func buildUpdateEmailURL(ctx context.Context, token string) (string, error) {
	u, err := url.Parse(ctxutils.HTTPReferer(ctx))
	if err != nil {
		return "", nil
	}

	u.Path = path.Join(u.Path, "users", "email", "update", "confirm")

	q := u.Query()
	q.Add("token", token)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func buildInvitationURL(ctx context.Context, token, email string, isUserExists bool) (string, error) {
	u, err := url.Parse(ctxutils.HTTPReferer(ctx))
	if err != nil {
		return "", nil
	}

	u.Path = path.Join(u.Path, "users", "invitation", "activate")

	q := u.Query()
	q.Add("token", token)
	q.Add("email", email)
	q.Add("isUserExists", strconv.FormatBool(isUserExists))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func buildSaveAuthURL(subdomain string) string {
	if subdomain == "" {
		return ""
	}
	if config.Config.Env == config.EnvLocal {
		return fmt.Sprintf("http://%s.%s:8080%s", subdomain, config.Config.Domain, model.SaveAuthPath)
	}
	return fmt.Sprintf("https://%s.%s%s", subdomain, config.Config.Domain, model.SaveAuthPath)
}

func buildServiceURL(subdomain string) string {
	if subdomain == "" {
		return ""
	}
	if config.Config.Env == config.EnvLocal {
		return fmt.Sprintf("http://%s.%s:5173", subdomain, config.Config.Domain)
	}
	return fmt.Sprintf("https://%s.%s", subdomain, config.Config.Domain)
}

func buildServiceDomain(subdomain string) string {
	if subdomain == "" {
		return ""
	}
	if config.Config.Env == config.EnvLocal {
		return fmt.Sprintf("%s.%s:5173", subdomain, config.Config.Domain)
	}
	return fmt.Sprintf("%s.%s", subdomain, config.Config.Domain)
}
