package internal

import (
	"errors"
	"strings"
)

func GetSubdomainFromHost(host string) (string, error) {
	if host == "" {
		return "", errors.New("empty host")
	}
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid host format")
	}
	return parts[0], nil
}
