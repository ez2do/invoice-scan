package httputil

import (
	"fmt"
	"invoice-scan/backend/pkg"
	"invoice-scan/backend/pkg/log"

	"net/http"
)

type HealthCheckResponse struct {
	Message string                    `json:"message,omitempty"`
	Status  pkg.ServiceHealthState    `json:"status"`
	Errors  []*pkg.ServiceHealthError `json:"errors,omitempty"`
}

type healthCheckHandler struct {
	checkFunctions []pkg.ServiceHealth
}

// NewHealthCheck return standard health check handler with list of pkg.ServiceHealth services
func NewHealthCheck(serviceHealths ...pkg.ServiceHealth) http.Handler {
	return &healthCheckHandler{checkFunctions: serviceHealths}
}

func (h *healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			log.Errorw("Panic when health check", "error", p)
			RespondError(w, http.StatusInternalServerError, fmt.Sprintf("health check fail: %v", p))
			return
		}
	}()
	var errs []*pkg.ServiceHealthError
	for _, f := range h.checkFunctions {
		if err := f.Check(); err != nil {
			errs = append(errs, &pkg.ServiceHealthError{
				Service: f.Name(),
				Error:   err.Error(),
			})
		}
	}

	if len(errs) == 0 {
		RespondJSON(w, http.StatusOK, HealthCheckResponse{
			Message: "Service is healthy",
			Status:  pkg.ServiceHealthy,
		})
		return
	}

	RespondJSON(w, http.StatusServiceUnavailable, HealthCheckResponse{
		Message: "Service is unhealthy",
		Errors:  errs,
		Status:  pkg.ServiceUnhealthy,
	})
}
