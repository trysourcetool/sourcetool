package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/auth"
)

type CookieConfig struct {
	tmpAuthDomain string
	isLocalEnv    bool
}

func NewCookieConfig() *CookieConfig {
	return &CookieConfig{
		tmpAuthDomain: config.Config.AuthDomain(),
		isLocalEnv:    config.Config.Env == config.EnvLocal,
	}
}

func (c *CookieConfig) getXSRFTokenSameSite() http.SameSite {
	if c.isLocalEnv {
		return http.SameSiteLaxMode
	}
	return http.SameSiteNoneMode
}

func (c *CookieConfig) isSecure() bool {
	return !c.isLocalEnv
}

func (c *CookieConfig) setCookie(w http.ResponseWriter, name, value string, maxAge int, httpOnly bool, sameSite http.SameSite, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: httpOnly,
		Secure:   c.isSecure(),
		SameSite: sameSite,
	})
}

func (c *CookieConfig) deleteCookie(w http.ResponseWriter, r *http.Request, name string, httpOnly bool, sameSite http.SameSite, domain string) {
	if cookie, _ := r.Cookie(name); cookie != nil {
		cookie.MaxAge = -1
		cookie.Domain = domain
		cookie.Path = "/"
		cookie.HttpOnly = httpOnly
		cookie.Secure = c.isSecure()
		cookie.SameSite = sameSite
		http.SetCookie(w, cookie)
	}
}

func (c *CookieConfig) SetTmpAuthCookie(w http.ResponseWriter, token, xsrfToken, domain string) {
	maxAge := int(auth.TmpTokenExpiration.Seconds())
	xsrfTokenSameSite := c.getXSRFTokenSameSite()

	c.setCookie(w, "access_token", token, maxAge, true, http.SameSiteStrictMode, domain)
	c.setCookie(w, "xsrf_token", xsrfToken, maxAge, false, xsrfTokenSameSite, domain)
	c.setCookie(w, "xsrf_token_same_site", xsrfToken, maxAge, true, http.SameSiteStrictMode, domain)
}

func (c *CookieConfig) DeleteTmpAuthCookie(w http.ResponseWriter, r *http.Request) {
	xsrfTokenSameSite := c.getXSRFTokenSameSite()

	c.deleteCookie(w, r, "access_token", true, http.SameSiteStrictMode, c.tmpAuthDomain)
	c.deleteCookie(w, r, "xsrf_token", false, xsrfTokenSameSite, c.tmpAuthDomain)
	c.deleteCookie(w, r, "xsrf_token_same_site", true, http.SameSiteStrictMode, c.tmpAuthDomain)
}

func (c *CookieConfig) SetAuthCookie(w http.ResponseWriter, token, refreshToken, xsrfToken string, tokenMaxAge, refreshTokenMaxAge, xsrfTokenMaxAge int, domain string) {
	xsrfTokenSameSite := c.getXSRFTokenSameSite()

	c.setCookie(w, "access_token", token, tokenMaxAge, true, http.SameSiteStrictMode, domain)
	c.setCookie(w, "refresh_token", refreshToken, refreshTokenMaxAge, true, http.SameSiteStrictMode, domain)
	c.setCookie(w, "xsrf_token", xsrfToken, xsrfTokenMaxAge, false, xsrfTokenSameSite, domain)
	c.setCookie(w, "xsrf_token_same_site", xsrfToken, xsrfTokenMaxAge, true, http.SameSiteStrictMode, domain)
}

func (c *CookieConfig) DeleteAuthCookie(w http.ResponseWriter, r *http.Request, domain string) {
	xsrfTokenSameSite := c.getXSRFTokenSameSite()

	c.deleteCookie(w, r, "access_token", true, http.SameSiteStrictMode, domain)
	c.deleteCookie(w, r, "refresh_token", true, http.SameSiteStrictMode, domain)
	c.deleteCookie(w, r, "xsrf_token", false, xsrfTokenSameSite, domain)
	c.deleteCookie(w, r, "xsrf_token_same_site", true, http.SameSiteStrictMode, domain)
}
