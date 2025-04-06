package adapters

import (
	"fmt"
	"time"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

func HealthDTOToResponse(health *dto.Health) *responses.HealthResponse {
	if health == nil {
		return nil
	}

	details := make(map[string]string)
	for k, v := range health.Details {
		details[k] = string(v)
	}

	return &responses.HealthResponse{
		Status:    string(health.Status),
		Version:   health.Version,
		Uptime:    formatUptime(health.Uptime),
		Timestamp: health.Timestamp.Format(time.RFC3339),
		Details:   details,
	}
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
