package auth

import (
	"path"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

func buildLoginURL(subdomain string) (string, error) {
	return internal.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("login"), nil)
}

func buildMagicLinkURL(subdomain, token string) (string, error) {
	base := config.Config.AuthBaseURL()
	if subdomain != "" && subdomain != "auth" {
		base = config.Config.OrgBaseURL(subdomain)
	}
	return internal.BuildURL(base, path.Join("auth", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}

func buildInvitationMagicLinkURL(subdomain, token string) (string, error) {
	baseURL := config.Config.OrgBaseURL(subdomain)
	return internal.BuildURL(baseURL, path.Join("auth", "invitations", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}
