package health

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/infra"
)

type StoreCE struct {
	db infra.DB
}

func NewStoreCE(db infra.DB) *StoreCE {
	return &StoreCE{
		db: db,
	}
}

func (s *StoreCE) Ping(ctx context.Context) (map[string]dto.HealthStatus, error) {
	details := make(map[string]dto.HealthStatus)

	details["postgres"] = s.checkPostgres(ctx)
	
	details["redis"] = s.checkRedis(ctx)

	return details, nil
}

func (s *StoreCE) checkPostgres(ctx context.Context) dto.HealthStatus {
	if db, ok := s.db.(interface{ DB() *sqlx.DB }); ok {
		if err := db.DB().PingContext(ctx); err != nil {
			return dto.HealthStatusDown
		}
		return dto.HealthStatusUp
	}
	return dto.HealthStatusDown
}

func (s *StoreCE) checkRedis(ctx context.Context) dto.HealthStatus {
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
