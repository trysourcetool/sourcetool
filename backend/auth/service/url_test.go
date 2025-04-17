package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trysourcetool/sourcetool/backend/auth"
	"github.com/trysourcetool/sourcetool/backend/config"
)

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
			expected:  "https://example.com" + auth.SaveAuthPath,
		},
		{
			name:      "empty subdomain",
			subdomain: "",
			expected:  "https://example.com" + auth.SaveAuthPath,
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
