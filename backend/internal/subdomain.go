package internal

import (
	"errors"
	"strings"
)

func GetSubdomainFromHost(host string) (string, error) {
	if host == "" {
		return "", errors.New("empty host")
	}

	// Strip port if present
	host = strings.Split(host, ":")[0]

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid host format")
	}

	// No subdomain if only domain and tld are present
	if len(parts) == 2 {
		return "", nil
	}

	return parts[0], nil
}
