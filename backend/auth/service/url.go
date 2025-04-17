package service

import (
	"path"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/utils/urlutil"
)

func buildLoginURL(subdomain string) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("login"), nil)
}

func buildMagicLinkURL(subdomain, token string) (string, error) {
	base := config.Config.AuthBaseURL()
	if subdomain != "" && subdomain != "auth" {
		base = config.Config.OrgBaseURL(subdomain)
	}
	return urlutil.BuildURL(base, path.Join("auth", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}

func buildInvitationMagicLinkURL(subdomain, token string) (string, error) {
	baseURL := config.Config.OrgBaseURL(subdomain)
	return urlutil.BuildURL(baseURL, path.Join("auth", "invitations", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}
