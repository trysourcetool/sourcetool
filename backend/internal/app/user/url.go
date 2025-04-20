package user

import (
	"path"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

func buildUpdateEmailURL(subdomain, token string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "email", "update", "confirm"), map[string]string{
		"token": token,
	})
}

func buildInvitationURL(subdomain, token, email string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("auth", "invitations", "login"), map[string]string{
		"token": token,
		"email": email,
	})
}
