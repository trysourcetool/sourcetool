package user

import (
	"path"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/utils/urlutil"
)

func buildUserActivateURL(token string) (string, error) {
	return urlutil.BuildURL(config.Config.AuthBaseURL(), path.Join("signup", "activate"), map[string]string{
		"token": token,
	})
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

func buildUpdateEmailURL(subdomain, token string) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "email", "update", "confirm"), map[string]string{
		"token": token,
	})
}

func buildInvitationURL(subdomain, token, email string, isUserExists bool) (string, error) {
	return urlutil.BuildURL(config.Config.OrgBaseURL(subdomain), path.Join("users", "invitation", "activate"), map[string]string{
		"token":        token,
		"email":        email,
		"isUserExists": strconv.FormatBool(isUserExists),
	})
}

func buildSaveAuthURL(subdomain string) (string, error) {
	return config.Config.OrgBaseURL(subdomain) + model.SaveAuthPath, nil
}
