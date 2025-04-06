package health

import "time"

type Status string

const (
	StatusUp Status = "up"
	StatusDown Status = "down"
)

type HealthDTO struct {
	Status    Status
	Version   string
	Uptime    time.Duration
	Timestamp time.Time
	Details   map[string]Status
}
