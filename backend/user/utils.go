package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/httputils"
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

func buildUserActivateURL(subdomain, token string) (string, error) {
	serviceURL, err := buildServiceURL(subdomain)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(serviceURL)
	if err != nil {
		return "", nil
	}

	u.Path = path.Join(u.Path, "signup", "activate")

	q := u.Query()
	q.Add("token", token)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func buildUpdateEmailURL(subdomain, token string) (string, error) {
	serviceURL, err := buildServiceURL(subdomain)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(serviceURL)
	if err != nil {
		return "", nil
	}

	u.Path = path.Join(u.Path, "users", "email", "update", "confirm")

	q := u.Query()
	q.Add("token", token)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func buildInvitationURL(subdomain, token, email string, isUserExists bool) (string, error) {
	serviceURL, err := buildServiceURL(subdomain)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(serviceURL)
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

type URLConfig struct {
	Port      string
	Path      string
	IsCloud   bool
	Subdomain string
	Env       string
	BaseURL   string
}

func (cfg URLConfig) buildURL() (string, error) {
	if cfg.IsCloud && cfg.Subdomain == "" {
		return "", fmt.Errorf("subdomain is required for cloud edition")
	}

	if cfg.Env == config.EnvLocal {
		if cfg.IsCloud {
			return fmt.Sprintf("http://%s.local.trysourcetool.com:%s%s",
				cfg.Subdomain, cfg.Port, cfg.Path), nil
		}
		return fmt.Sprintf("http://localhost:%s%s", cfg.Port, cfg.Path), nil
	}

	if cfg.IsCloud {
		domain, err := httputils.ExtractDomainFromURL(cfg.BaseURL)
		if err != nil {
			return "", err
		}
		scheme := strings.Split(cfg.BaseURL, "://")[0]
		return fmt.Sprintf("%s://%s.%s%s", scheme, cfg.Subdomain, domain, cfg.Path), nil
	}

	if cfg.Path != "" {
		return cfg.BaseURL + cfg.Path, nil
	}
	return cfg.BaseURL, nil
}

func buildSaveAuthURL(subdomain string) (string, error) {
	cfg := URLConfig{
		Port:      "8080",
		Path:      model.SaveAuthPath,
		IsCloud:   config.Config.IsCloudEdition,
		Subdomain: subdomain,
		Env:       config.Config.Env,
		BaseURL:   config.Config.BaseURL,
	}
	return cfg.buildURL()
}

func buildServiceURL(subdomain string) (string, error) {
	cfg := URLConfig{
		Port:      "5173",
		Path:      "",
		IsCloud:   config.Config.IsCloudEdition,
		Subdomain: subdomain,
		Env:       config.Config.Env,
		BaseURL:   config.Config.BaseURL,
	}
	return cfg.buildURL()
}

func buildServiceDomain(subdomain string) (string, error) {
	if config.Config.Env == config.EnvLocal {
		if config.Config.IsCloudEdition {
			return fmt.Sprintf("%s.local.trysourcetool.com", subdomain), nil
		}
		return "localhost", nil
	}

	if config.Config.IsCloudEdition {
		domain, err := httputils.ExtractDomainFromURL(config.Config.BaseURL)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.%s", subdomain, domain), nil
	}

	domain, err := httputils.ExtractDomainFromURL(config.Config.BaseURL)
	if err != nil {
		return "", err
	}
	return domain, nil
}
