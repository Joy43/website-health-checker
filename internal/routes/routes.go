package routes

import (
	"log/slog"
	"net/http"

	"github.com/ssjoy/website-health-checker/internal/handlers"
	"github.com/ssjoy/website-health-checker/internal/middleware"
)

// Setup creates the main HTTP handler with routing and middlewares.
func Setup(logger *slog.Logger, hh *handlers.HealthHandler, ch *handlers.CacheHandler) http.Handler {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/health", hh.GetHealth)
	mux.HandleFunc("/health/details", hh.GetDetailedHealth)
	mux.HandleFunc("/cache", ch.HandleCache)

	// Middlewares chain: Request -> Logger -> Recovery -> Mux
	// Recovery is inner to Logger, so that panics are caught and the 500 response is logged.
	var handler http.Handler = mux
	handler = middleware.Recovery(logger)(handler)
	handler = middleware.Logger(logger)(handler)

	return handler
}
