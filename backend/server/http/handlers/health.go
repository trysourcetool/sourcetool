package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/health"
	"github.com/trysourcetool/sourcetool/backend/utils/httputil"
)

type HealthHandler struct {
	service health.Service
}

func NewHealthHandler(service health.Service) *HealthHandler {
	return &HealthHandler{service}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	healthStatus, err := h.service.Check(r.Context())
	if err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}

	statusCode := http.StatusOK
	if healthStatus.Status == health.StatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	if err := httputil.WriteJSON(w, statusCode, healthStatus); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
