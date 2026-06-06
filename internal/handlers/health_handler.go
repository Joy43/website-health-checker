package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ssjoy/website-health-checker/internal/services"
)

// HealthHandler handles health-related HTTP endpoints.
type HealthHandler struct {
	service services.HealthService
}

// NewHealthHandler returns a new HealthHandler instance.
func NewHealthHandler(service services.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// GetHealth handles GET /health requests.
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"Method Not Allowed"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	status := h.service.GetHealth(r.Context())
	
	// If any services are disconnected, we should return a 503 Service Unavailable or just 200 OK.
	// The prompt's expected response format:
	// { "status": "ok", "service": "health-checker", ... }
	// Let's return 200 OK, but if status is "error" we could return 503 Service Unavailable as a production standard.
	// Let's set StatusServiceUnavailable if status.Status == "error" or just return 200 OK for standard health check.
	// Actually, Kubernetes and other health checkers expect non-200 codes when failing, so setting status to 503 is best practice!
	if status.Status == "error" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(status)
}

// GetDetailedHealth handles GET /health/details requests.
func (h *HealthHandler) GetDetailedHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"Method Not Allowed"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	status := h.service.GetDetailedHealth(r.Context())

	if status.MySQL.Status == "unhealthy" || status.Redis.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(status)
}
