package handlers

import (
	"net/http"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/health"
	"github.com/trysourcetool/sourcetool/backend/server/http/adapters"
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

	healthResponse := adapters.HealthDTOToResponse(healthDTO)

	statusCode := http.StatusOK
	if healthDTO.Status == dto.HealthStatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	if err := httputil.WriteJSON(w, statusCode, healthResponse); err != nil {
		httputil.WriteErrJSON(r.Context(), w, err)
		return
	}
}
