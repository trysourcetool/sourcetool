//go:build !ee

package license

import (
	"context"
	"time"
)

func NewChecker(baseURL, licenseKey string, timeout time.Duration) (*Checker, error) {
	return &Checker{}, nil
}

func (c *Checker) Check(ctx context.Context) error {
	return nil
}
