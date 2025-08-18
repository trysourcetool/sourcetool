//go:build !ee

package license

import (
	"context"
	"time"
)

func NewChecker(baseURL, licenseKey string, timeout time.Duration) (*Checker, error) {
	return &Checker{}, nil
}

func (c *Checker) Validate(ctx context.Context) error {
	return nil
}

func (c *Checker) UpdateSeats(ctx context.Context, seats int64) error {
	return nil
}
