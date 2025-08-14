//go:build ee

package license

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

var licenseServerURLRegex = regexp.MustCompile(`^https?://(?:[a-zA-Z0-9-]+\.)?license\.trysourcetool\.com$`)

type SubscriptionResponse struct {
	ID                   string `json:"id"`
	UserID               string `json:"userId"`
	PlanID               string `json:"planId"`
	Status               string `json:"status"`
	StripeCustomerID     string `json:"stripeCustomerId"`
	StripeSubscriptionID string `json:"stripeSubscriptionId"`
	TrialStart           string `json:"trialStart"`
	TrialEnd             string `json:"trialEnd"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
}

type LicenseValidityResponse struct {
	Valid        bool                  `json:"valid"`
	Status       string                `json:"status"`
	Subscription *SubscriptionResponse `json:"subscription,omitempty"`
}

type UpdateSeatsRequest struct {
	Seats int64 `json:"seats"`
}

func NewChecker(baseURL, licenseKey string, timeout time.Duration) (*Checker, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("license server base URL cannot be empty")
	}
	if config.Config.Env != config.EnvLocal {
		if !licenseServerURLRegex.MatchString(baseURL) {
			return nil, fmt.Errorf("license server URL must be a *.license.trysourcetool.com domain")
		}
	}
	return &Checker{BaseURL: baseURL, LicenseKey: licenseKey, Timeout: timeout}, nil
}

func (c *Checker) Validate(ctx context.Context) error {
	endpoint := c.BaseURL + "/v1/validate"

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Sourcetool-License-Key", c.LicenseKey)

	timeout := c.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("license server returned status: %s", resp.Status)
	}

	var result LicenseValidityResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if !result.Valid {
		return fmt.Errorf("license invalid: %s", result.Status)
	}
	return nil
}

func (c *Checker) UpdateSeats(ctx context.Context, seats int64) error {
	endpoint := c.BaseURL + "/v1/seats"

	reqBody := UpdateSeatsRequest{Seats: seats}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Sourcetool-License-Key", c.LicenseKey)
	req.Header.Set("Content-Type", "application/json")

	timeout := c.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("license server returned status: %s", resp.Status)
	}
	return nil
}
