package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a Gin middleware that logs incoming HTTP requests using slog.
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start)

		logger.Info("HTTP Request",
			slog.String("method", c.Request.Method),
			slog.String("uri", c.Request.RequestURI),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", duration),
		)
	}
}
