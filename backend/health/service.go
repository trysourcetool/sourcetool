package health

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type Service interface {
	Check(ctx context.Context) (*dto.Health, error)
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

func (s *ServiceCE) Check(ctx context.Context) (*dto.Health, error) {
	details := make(map[string]dto.HealthStatus)

	details["postgres"] = s.checkPostgres(ctx)
	details["redis"] = s.checkRedis(ctx)

	overallStatus := dto.HealthStatusUp
	for _, status := range details {
		if status == dto.HealthStatusDown {
			overallStatus = dto.HealthStatusDown
			break
		}
	}

	return &dto.Health{
		Status:    overallStatus,
		Version:   "1.0", // Using the version from Swagger docs
		Uptime:    time.Since(s.startTime),
		Timestamp: time.Now().UTC(),
		Details:   details,
	}, nil
}

func (s *ServiceCE) checkPostgres(ctx context.Context) dto.HealthStatus {
	if db, ok := s.Store.(interface{ DB() *sqlx.DB }); ok {
		if err := db.DB().PingContext(ctx); err != nil {
			return dto.HealthStatusDown
		}
		return dto.HealthStatusUp
	}
	return dto.HealthStatusDown
}

func (s *ServiceCE) checkRedis(ctx context.Context) dto.HealthStatus {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
	})
	defer client.Close()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return dto.HealthStatusDown
	}
	return dto.HealthStatusUp
}
