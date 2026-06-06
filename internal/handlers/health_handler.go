package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Summary      Get overall system health
// @Description  Checks API status, MySQL database connectivity, and Redis cache connectivity.
// @Tags         health
// @Produce      json
// @Success      200  {object}  services.HealthStatus
// @Failure      503  {object}  services.HealthStatus
// @Router       /health [get]
func (h *HealthHandler) GetHealth(c *gin.Context) {
	status := h.service.GetHealth(c.Request.Context())
	
	if status.Status == "error" {
		c.JSON(http.StatusServiceUnavailable, status)
	} else {
		c.JSON(http.StatusOK, status)
	}
}

// GetDetailedHealth handles GET /health/details requests.
// @Summary      Get detailed system health status
// @Description  Retrieves status details for MySQL, Redis, and overall system uptime.
// @Tags         health
// @Produce      json
// @Success      200  {object}  services.DetailedHealthStatus
// @Failure      503  {object}  services.DetailedHealthStatus
// @Router       /health/details [get]
func (h *HealthHandler) GetDetailedHealth(c *gin.Context) {
	status := h.service.GetDetailedHealth(c.Request.Context())

	if status.MySQL.Status == "unhealthy" || status.Redis.Status == "unhealthy" {
		c.JSON(http.StatusServiceUnavailable, status)
	} else {
		c.JSON(http.StatusOK, status)
	}
}
