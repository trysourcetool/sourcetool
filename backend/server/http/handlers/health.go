package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/health"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type HealthHandler struct {
	service health.Service
}

func NewHealthHandler(service health.Service) *HealthHandler {
	return &HealthHandler{service}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	healthDTO, err := h.service.Check(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	healthResponse := toHealthResponse(healthDTO)

	statusCode := http.StatusOK
	if healthDTO.Status == dto.HealthStatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	if err := httputil.WriteJSON(w, statusCode, healthResponse); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}

func toHealthResponse(dto *dto.Health) *responses.HealthResponse {
	details := make(map[string]string)
	for k, v := range dto.Details {
		details[k] = string(v)
	}

	return &responses.HealthResponse{
		Status:    string(dto.Status),
		Version:   dto.Version,
		Uptime:    formatUptime(dto.Uptime),
		Timestamp: dto.Timestamp.Format(time.RFC3339),
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
