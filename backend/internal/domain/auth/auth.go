package auth

import (
	"time"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

const (
	EmailTokenExpiration     = time.Duration(24) * time.Hour
	tokenExpiration          = time.Duration(60) * time.Minute
	tokenExpirationDev       = time.Duration(365*24) * time.Hour
	RefreshTokenExpiration   = time.Duration(30*24) * time.Hour
	XSRFTokenExpiration      = time.Duration(30*24) * time.Hour
	RefreshTokenMaxAgeBuffer = time.Duration(7*24) * time.Hour
	TmpTokenExpiration       = time.Duration(30) * time.Minute

	SaveAuthPath = "/api/v1/auth/save"
)

func TokenExpiration() time.Duration {
	if config.Config.Env == config.EnvLocal {
		return tokenExpirationDev
	}
	return tokenExpiration
}
