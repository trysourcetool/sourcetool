package config

import (
	"log"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v11"
)

var Config *config

const (
	EnvLocal   = "local"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

type config struct {
	Env            string `env:"ENV"`
	BaseURL        string `env:"BASE_URL"`
	SSL            bool   `env:"-"`
	Protocol       string `env:"-"`
	BaseDomain     string `env:"-"`
	BaseHostname   string `env:"-"`
	IsCloudEdition bool   `env:"-"`
	EncryptionKey  string `env:"ENCRYPTION_KEY"`
	Jwt            struct {
		Key string `env:"JWT_KEY"`
	}
	Postgres struct {
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DB       string `env:"POSTGRES_DB"`
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST"`
		Port     string `env:"REDIS_PORT"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
	}
	Google struct {
		OAuth struct {
			// Optional for self-hosted deployments, users may leave these unset if not using Google OAuth.
			ClientID     string `env:"GOOGLE_OAUTH_CLIENT_ID" envDefault:""`
			ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET" envDefault:""`
		}
	}
	SMTP struct {
		Host      string `env:"SMTP_HOST"`
		Port      string `env:"SMTP_PORT"`
		Username  string `env:"SMTP_USERNAME"`
		Password  string `env:"SMTP_PASSWORD"`
		FromEmail string `env:"SMTP_FROM_EMAIL"`
		UseTLS    bool   `env:"SMTP_USE_TLS" envDefault:"true"`
	}
}

func Init() {
	cfg := new(config)
	envOpts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, envOpts); err != nil {
		log.Fatal("[INIT] config: ", err)
	}

	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	cfg.IsCloudEdition = isCloudEdition(cfg.BaseURL)

	baseURLParts := strings.Split(cfg.BaseURL, "://")
	if len(baseURLParts) != 2 {
		log.Fatal("[INIT] invalid BASE_URL format: ", cfg.BaseURL)
	}
	cfg.Protocol = baseURLParts[0]
	cfg.BaseHostname = baseURLParts[1]
	cfg.SSL = cfg.Protocol == "https"

	hostnameParts := strings.Split(cfg.BaseHostname, ":")
	cfg.BaseDomain = hostnameParts[0]

	Config = cfg
}

// AuthHostname returns the hostname for the Auth domain.
func (c *config) AuthHostname() string {
	if c.IsCloudEdition {
		return "auth." + c.BaseHostname
	}
	return c.BaseHostname
}

// OrgHostname returns the hostname for the organization.
func (c *config) OrgHostname(subdomain string) string {
	if c.IsCloudEdition {
		return subdomain + "." + c.BaseHostname
	}
	return c.BaseHostname
}

// AuthDomain returns the domain for the Auth domain.
func (c *config) AuthDomain() string {
	if c.IsCloudEdition {
		return "auth." + c.BaseDomain
	}
	return c.BaseDomain
}

// OrgDomain returns the domain for the organization.
func (c *config) OrgDomain(subdomain string) string {
	if c.IsCloudEdition {
		return subdomain + "." + c.BaseDomain
	}
	return c.BaseDomain
}

// AuthBaseURL returns the base URL for the Auth domain.
func (c *config) AuthBaseURL() string {
	return c.Protocol + "://" + c.AuthHostname()
}

// OrgBaseURL returns the base URL for the organization domain.
func (c *config) OrgBaseURL(subdomain string) string {
	return c.Protocol + "://" + c.OrgHostname(subdomain)
}

// WebSocketOrgBaseURL returns the base URL for the organization domain for WebSocket connections.
func (c *config) WebSocketOrgBaseURL(subdomain string) string {
	if c.SSL {
		return "wss://" + c.OrgHostname(subdomain)
	}
	return "ws://" + c.OrgHostname(subdomain)
}

func isCloudEdition(baseURL string) bool {
	cloudDomainRegex := regexp.MustCompile(`^https?://(?:([^.]+)\.)?trysourcetool\.com(?::\d+)?$`)
	matches := cloudDomainRegex.FindStringSubmatch(baseURL)
	return len(matches) > 1
}
