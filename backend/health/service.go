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
	details, err := s.Store.Health().Ping(ctx)
	if err != nil {
		return nil, err
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
