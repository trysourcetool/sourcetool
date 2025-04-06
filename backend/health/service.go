package health

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
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
	pgDetails, err := s.Store.Health().Ping(ctx)
	if err != nil {
		return nil, err
	}
	
	details := make(map[string]model.HealthStatus)
	for k, v := range pgDetails {
		details[k] = v
	}
	
	if err := s.Memory.Redis().Ping(ctx); err != nil {
		details["redis"] = model.HealthStatusDown
	} else {
		details["redis"] = model.HealthStatusUp
	}

	overallStatus := model.HealthStatusUp
	for _, status := range details {
		if status == model.HealthStatusDown {
			overallStatus = model.HealthStatusDown
			break
		}
	}

	dtoDetails := make(map[string]dto.HealthStatus)
	for k, v := range details {
		if v == model.HealthStatusUp {
			dtoDetails[k] = dto.HealthStatusUp
		} else {
			dtoDetails[k] = dto.HealthStatusDown
		}
	}

	dtoStatus := dto.HealthStatusUp
	if overallStatus == model.HealthStatusDown {
		dtoStatus = dto.HealthStatusDown
	}

	return &dto.Health{
		Status:    dtoStatus,
		Uptime:    time.Since(s.startTime),
		Timestamp: time.Now().UTC(),
		Details:   dtoDetails,
	}, nil
}
