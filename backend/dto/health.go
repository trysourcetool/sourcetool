package dto

import "time"

type HealthStatus string

const (
	HealthStatusUp HealthStatus = "up"
	HealthStatusDown HealthStatus = "down"
)

type Health struct {
	Status    HealthStatus
	Version   string
	Uptime    time.Duration
	Timestamp time.Time
	Details   map[string]HealthStatus
}
