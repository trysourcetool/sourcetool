package license

import "time"

type Checker struct {
	BaseURL    string
	LicenseKey string
	Timeout    time.Duration
}
