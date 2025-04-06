package dto

import "time"

type HealthStatus string

const (
	HealthStatusUp HealthStatus = "up"
	HealthStatusDown HealthStatus = "down"
)

type Health struct {
	Status    HealthStatus
	Uptime    time.Duration
	Timestamp time.Time
	Details   map[string]HealthStatus
}
