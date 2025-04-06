package health

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

type Service interface {
	Check(ctx context.Context) (*HealthStatus, error)
}

type HealthStatus struct {
	Status    Status            `json:"status"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Timestamp string            `json:"timestamp"`
	Details   map[string]Status `json:"details"`
}

type ServiceCE struct {
	*infra.Dependency
	startTime time.Time
}

func NewServiceCE(dep *infra.Dependency) *ServiceCE {
	return &ServiceCE{
		Dependency: dep,
		startTime:  time.Now(),
	}
}

func (s *ServiceCE) Check(ctx context.Context) (*HealthStatus, error) {
	details := make(map[string]Status)

	details["postgres"] = s.checkPostgres(ctx)

	details["redis"] = s.checkRedis(ctx)

	overallStatus := StatusUp
	for _, status := range details {
		if status == StatusDown {
			overallStatus = StatusDown
			break
		}
	}

	return &HealthStatus{
		Status:    overallStatus,
		Version:   "1.0", // Using the version from Swagger docs
		Uptime:    formatUptime(time.Since(s.startTime)),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details:   details,
	}, nil
}

func (s *ServiceCE) checkPostgres(ctx context.Context) Status {
	if db, ok := s.Store.(interface{ DB() *sqlx.DB }); ok {
		if err := db.DB().PingContext(ctx); err != nil {
			return StatusDown
		}
		return StatusUp
	}
	return StatusDown
}

func (s *ServiceCE) checkRedis(ctx context.Context) Status {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
	})
	defer client.Close()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return StatusDown
	}
	return StatusUp
}

func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
