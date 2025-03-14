package config

import (
	"log"
	"regexp"

	"github.com/caarlos0/env/v11"
)

var Config *ConfigCE

const (
	EnvLocal   = "local"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

type ConfigCE struct {
	Env            string
	IsCloudEdition bool
	Domain         string `env:"DOMAIN"`
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
		Password string `env:"REDIS_PASSWORD"`
	}
	Google struct {
		OAuth struct {
			ClientID     string `env:"GOOGLE_OAUTH_CLIENT_ID"`
			ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET"`
			CallbackURL  string `env:"GOOGLE_OAUTH_CALLBACK_URL"`
		}
	}
	SMTP struct {
		Host      string `env:"SMTP_HOST"`
		Port      string `env:"SMTP_PORT"`
		Username  string `env:"SMTP_USERNAME"`
		Password  string `env:"SMTP_PASSWORD"`
		FromEmail string `env:"SMTP_FROM_EMAIL"`
	}
}

func Init() {
	cfg := new(ConfigCE)
	envOpts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, envOpts); err != nil {
		log.Fatal("[INIT] config: ", err)
	}

	localhostRegex := regexp.MustCompile(`^localhost(:\d+)?$`)
	cloudDomainRegex := regexp.MustCompile(`^(?:([^.]+)\.)?trysourcetool\.com$`)

	// Set Env and IsCloudEdition based on DOMAIN
	if localhostRegex.MatchString(cfg.Domain) {
		cfg.Env = EnvLocal
		cfg.IsCloudEdition = false
	} else if matches := cloudDomainRegex.FindStringSubmatch(cfg.Domain); matches != nil {
		cfg.IsCloudEdition = true
		if matches[1] != "" {
			cfg.Env = matches[1] // subdomain exists (e.g., staging.trysourcetool.com)
		} else {
			cfg.Env = EnvProd // just trysourcetool.com
		}
	} else {
		cfg.Env = EnvProd
		cfg.IsCloudEdition = false
	}

	log.Printf("env: %s, isCloudEdition: %t", cfg.Env, cfg.IsCloudEdition)

	Config = cfg
}
