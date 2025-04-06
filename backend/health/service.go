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

	if db, ok := s.Store.(interface{ DB() *sqlx.DB }); ok {
		if err := db.DB().PingContext(ctx); err != nil {
			details["postgres"] = dto.HealthStatusDown
		} else {
			details["postgres"] = dto.HealthStatusUp
		}
	} else {
		details["postgres"] = dto.HealthStatusDown
	}

	if err := s.Memory.Redis().Ping(ctx); err != nil {
		details["redis"] = dto.HealthStatusDown
	} else {
		details["redis"] = dto.HealthStatusUp
	}

	overallStatus := dto.HealthStatusUp
	for _, status := range details {
		if status == dto.HealthStatusDown {
			overallStatus = dto.HealthStatusDown
			break
		}
	}

	return &dto.Health{
		Status:    overallStatus,
		Uptime:    time.Since(s.startTime),
		Timestamp: time.Now().UTC(),
		Details:   details,
	}, nil
}
