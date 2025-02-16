package user

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"path"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/model"
)

func generateSecret() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 32)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
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
