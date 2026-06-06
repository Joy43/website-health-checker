package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriterWrapper wraps standard http.ResponseWriter to capture status code.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	// If WriteHeader was not called, status is 200 OK
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Logger returns a middleware that logs incoming HTTP requests using slog.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default status code
			}

			next.ServeHTTP(wrapper, r)

			duration := time.Since(start)

			logger.Info("HTTP Request",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.String("ip", r.RemoteAddr),
				slog.Int("status", wrapper.statusCode),
				slog.Duration("duration", duration),
			)
		})
	}
}
