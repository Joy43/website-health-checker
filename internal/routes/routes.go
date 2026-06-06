package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ssjoy/website-health-checker/internal/handlers"
	"github.com/ssjoy/website-health-checker/internal/middleware"
)

// Setup creates the main Gin engine with routing and middlewares.
func Setup(logger *slog.Logger, hh *handlers.HealthHandler, ch *handlers.CacheHandler) *gin.Engine {
	// Set Gin mode to ReleaseMode to disable debugging logs in production
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global Middlewares
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))

	// Swagger Documentation UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Application Endpoints
	r.GET("/health", hh.GetHealth)
	r.GET("/health/details", hh.GetDetailedHealth)
	r.GET("/cache", ch.GetCache)
	r.POST("/cache", ch.SetCache)

	return r
}
