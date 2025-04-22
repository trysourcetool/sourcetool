package config

import (
	"strings"
	"testing"
)

func TestConfig_HostnameAndDomain(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		isCloudEdition bool
		subdomain      string
		wantAuthHost   string
		wantOrgHost    string
		wantAuthDomain string
		wantOrgDomain  string
	}{
		{
			name:           "local environment",
			baseURL:        "http://localhost:3000",
			isCloudEdition: false,
			subdomain:      "test",
			wantAuthHost:   "localhost:3000",
			wantOrgHost:    "localhost:3000",
			wantAuthDomain: "localhost",
			wantOrgDomain:  "localhost",
		},
		{
			name:           "non-cloud edition with custom domain",
			baseURL:        "https://custom-domain.com:8080",
			isCloudEdition: false,
			subdomain:      "test",
			wantAuthHost:   "custom-domain.com:8080",
			wantOrgHost:    "custom-domain.com:8080",
			wantAuthDomain: "custom-domain.com",
			wantOrgDomain:  "custom-domain.com",
		},
		{
			name:           "cloud edition",
			baseURL:        "https://example.com",
			isCloudEdition: true,
			subdomain:      "test",
			wantAuthHost:   "auth.example.com",
			wantOrgHost:    "test.example.com",
			wantAuthDomain: "auth.example.com",
			wantOrgDomain:  "test.example.com",
		},
		{
			name:           "cloud edition with port",
			baseURL:        "https://example.com:8080",
			isCloudEdition: true,
			subdomain:      "test",
			wantAuthHost:   "auth.example.com:8080",
			wantOrgHost:    "test.example.com:8080",
			wantAuthDomain: "auth.example.com",
			wantOrgDomain:  "test.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config{
				BaseURL:        tt.baseURL,
				IsCloudEdition: tt.isCloudEdition,
			}

			// Parse BaseURL to set Protocol and BaseHostname
			baseURLParts := strings.Split(tt.baseURL, "://")
			if len(baseURLParts) != 2 {
				t.Fatalf("invalid BASE_URL format: %s", tt.baseURL)
			}
			cfg.Protocol = baseURLParts[0]
			cfg.BaseHostname = baseURLParts[1]
			cfg.SSL = cfg.Protocol == "https"

			hostnameParts := strings.Split(cfg.BaseHostname, ":")
			cfg.BaseDomain = hostnameParts[0]

			// Test AuthHostname
			if got := cfg.AuthHostname(); got != tt.wantAuthHost {
				t.Errorf("AuthHostname() = %v, want %v", got, tt.wantAuthHost)
			}

			// Test OrgHostname
			if got := cfg.OrgHostname(tt.subdomain); got != tt.wantOrgHost {
				t.Errorf("OrgHostname() = %v, want %v", got, tt.wantOrgHost)
			}

			// Test AuthDomain
			if got := cfg.AuthDomain(); got != tt.wantAuthDomain {
				t.Errorf("AuthDomain() = %v, want %v", got, tt.wantAuthDomain)
			}

			// Test OrgDomain
			if got := cfg.OrgDomain(tt.subdomain); got != tt.wantOrgDomain {
				t.Errorf("OrgDomain() = %v, want %v", got, tt.wantOrgDomain)
			}

			// Test AuthBaseURL
			wantAuthBaseURL := cfg.Protocol + "://" + tt.wantAuthHost
			if got := cfg.AuthBaseURL(); got != wantAuthBaseURL {
				t.Errorf("AuthBaseURL() = %v, want %v", got, wantAuthBaseURL)
			}

			// Test OrgBaseURL
			wantOrgBaseURL := cfg.Protocol + "://" + tt.wantOrgHost
			if got := cfg.OrgBaseURL(tt.subdomain); got != wantOrgBaseURL {
				t.Errorf("OrgBaseURL() = %v, want %v", got, wantOrgBaseURL)
			}

			// Test WebSocketOrgBaseURL
			wantWSOrgBaseURL := "ws"
			if cfg.SSL {
				wantWSOrgBaseURL = "wss"
			}
			wantWSOrgBaseURL += "://" + tt.wantOrgHost
			if got := cfg.WebSocketOrgBaseURL(tt.subdomain); got != wantWSOrgBaseURL {
				t.Errorf("WebSocketOrgBaseURL() = %v, want %v", got, wantWSOrgBaseURL)
			}
		})
	}
}
