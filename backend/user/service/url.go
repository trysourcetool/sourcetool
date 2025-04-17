package service

import (
	"path"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/utils/urlutil"
)

func buildUpdateEmailURL(subdomain, token string) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "email", "update", "confirm"), map[string]string{
		"token": token,
	})
}

func buildInvitationURL(subdomain, token, email string) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("auth", "invitations", "login"), map[string]string{
		"token": token,
		"email": email,
	})
}
