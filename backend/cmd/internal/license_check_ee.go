//go:build ee

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

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

func CheckLicense() error {
	baseURL := os.Getenv("LICENSE_SERVER_BASE_URL")
	if baseURL == "" {
		baseURL = "http://host.docker.internal:8082"
	}
	endpoint := baseURL + "/v1/validate"
	licenseKey := os.Getenv("LICENSE_KEY")
	if licenseKey == "" {
		return errors.New("LICENSE_KEY is not set")
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Sourcetool-License-Key", licenseKey)

	client := &http.Client{Timeout: 10 * time.Second}
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
