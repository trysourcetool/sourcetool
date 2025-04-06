package health

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/trysourcetool/sourcetool/backend/config"
)

type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

type Service interface {
	Check(ctx context.Context) (*HealthStatus, error)
}

type ServiceCE struct {
	db        *sqlx.DB
	redisOpts *redis.Options
	startTime time.Time
}

type HealthStatus struct {
	Status    Status            `json:"status"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Timestamp string            `json:"timestamp"`
	Details   map[string]Status `json:"details"`
}

func NewServiceCE(db *sqlx.DB) *ServiceCE {
	return &ServiceCE{
		db: db,
		redisOpts: &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
			Password: config.Config.Redis.Password,
		},
		startTime: time.Now(),
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
	if err := s.db.PingContext(ctx); err != nil {
		return StatusDown
	}
	return StatusUp
}

func (s *ServiceCE) checkRedis(ctx context.Context) Status {
	client := redis.NewClient(s.redisOpts)
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
