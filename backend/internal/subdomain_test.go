package internal

import "testing"

func TestGetSubdomainFromHost(t *testing.T) {
	tests := []struct {
		host    string
		want    string
		wantErr bool
	}{
		{"org.trysourcetool.com", "org", false},
		{"org.trysourcetool.com:8080", "org", false},
		{"trysourcetool.com", "", false},
		{"trysourcetool.com:8080", "", false},
		{"localhost", "", true},
	}

	for _, tt := range tests {
		got, err := GetSubdomainFromHost(tt.host)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetSubdomainFromHost(%q) error = %v, wantErr %v", tt.host, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("GetSubdomainFromHost(%q) = %q, want %q", tt.host, got, tt.want)
		}
	}
}
