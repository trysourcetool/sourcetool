package model

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/dto"
)

type HealthStore interface {
	Ping(ctx context.Context) (map[string]dto.HealthStatus, error)
}
