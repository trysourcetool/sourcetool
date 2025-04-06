package health

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
)

type StoreCE struct {
	db infra.DB
}

func NewStoreCE(db infra.DB) *StoreCE {
	return &StoreCE{
		db: db,
	}
}

func (s *StoreCE) Ping(ctx context.Context) (map[string]model.HealthStatus, error) {
	details := make(map[string]model.HealthStatus)

	details["postgres"] = s.checkPostgres(ctx)

	return details, nil
}

func (s *StoreCE) checkPostgres(ctx context.Context) model.HealthStatus {
	if db, ok := s.db.(interface{ DB() *sqlx.DB }); ok {
		if err := db.DB().PingContext(ctx); err != nil {
			return model.HealthStatusDown
		}
		return model.HealthStatusUp
	}
	return model.HealthStatusDown
}
