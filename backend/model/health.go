package model

import (
	"context"
)

type HealthStatus string

const (
	HealthStatusUp HealthStatus = "up"
	HealthStatusDown HealthStatus = "down"
)

type HealthStore interface {
	Ping(ctx context.Context) (map[string]HealthStatus, error)
}
