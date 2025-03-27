package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/model"
)

func TestBuildUserActivateURL(t *testing.T) {
	// Setup test config
	config.Config = &config.Cfg{
		BaseURL:      "https://example.com",
		SSL:          true,
		Protocol:     "https",
		BaseDomain:   "example.com",
		BaseHostname: "example.com",
	}

	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "valid token",
			token:    "test-token",
			expected: "https://example.com/signup/activate?token=test-token",
		},
		{
			name:     "empty token",
			token:    "",
			expected: "https://example.com/signup/activate?token=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildUserActivateURL(tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildUpdateEmailURL(t *testing.T) {
	// Setup test config
	config.Config = &config.Cfg{
		BaseURL:      "https://example.com",
		SSL:          true,
		Protocol:     "https",
		BaseDomain:   "example.com",
		BaseHostname: "example.com",
	}

	tests := []struct {
		name      string
		subdomain string
		token     string
		expected  string
	}{
		{
			name:      "valid subdomain and token",
			subdomain: "test",
			token:     "test-token",
			expected:  "https://example.com/users/email/update/confirm?token=test-token",
		},
		{
			name:      "empty subdomain",
			subdomain: "",
			token:     "test-token",
			expected:  "https://example.com/users/email/update/confirm?token=test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildUpdateEmailURL(tt.subdomain, tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildInvitationURL(t *testing.T) {
	// Setup test config
	config.Config = &config.Cfg{
		BaseURL:      "https://example.com",
		SSL:          true,
		Protocol:     "https",
		BaseDomain:   "example.com",
		BaseHostname: "example.com",
	}

	tests := []struct {
		name      string
		subdomain string
		token     string
		email     string
		expected  string
	}{
		{
			name:      "new user invitation",
			subdomain: "test",
			token:     "test-token",
			email:     "test@example.com",
			expected:  "https://example.com/users/invitation/activate?email=test%40example.com&token=test-token",
		},
		{
			name:      "existing user invitation",
			subdomain: "test",
			token:     "test-token",
			email:     "existing@example.com",
			expected:  "https://example.com/users/invitation/activate?email=existing%40example.com&token=test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildInvitationURL(tt.subdomain, tt.token, tt.email)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildSaveAuthURL(t *testing.T) {
	// Setup test config
	config.Config = &config.Cfg{
		BaseURL:      "https://example.com",
		SSL:          true,
		Protocol:     "https",
		BaseDomain:   "example.com",
		BaseHostname: "example.com",
	}

	tests := []struct {
		name      string
		subdomain string
		expected  string
	}{
		{
			name:      "valid subdomain",
			subdomain: "test",
			expected:  "https://example.com" + model.SaveAuthPath,
		},
		{
			name:      "empty subdomain",
			subdomain: "",
			expected:  "https://example.com" + model.SaveAuthPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildSaveAuthURL(tt.subdomain)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
