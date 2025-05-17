package internal

import (
	"errors"
	"net"
	"strings"
)

func GetSubdomainFromHost(host string) (string, error) {
	if host == "" {
		return "", errors.New("empty host")
	}

	// Strip port if present
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

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
