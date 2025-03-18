package urlutil

import (
	"net/url"
	"regexp"
)

var cloudDomainRegex = regexp.MustCompile(`^https?://(?:([^.]+)\.)?trysourcetool\.com(?::\d+)?$`)

func IsCloudEdition(baseURL string) bool {
	matches := cloudDomainRegex.FindStringSubmatch(baseURL)
	return len(matches) > 1
}

func BuildURL(baseURL, path string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	u.Path = path

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
