package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery recovers from panics and logs stack traces.
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.Error("panic recovered",
						slog.Any("error", err),
						slog.String("stack", string(stack)),
						slog.String("uri", r.RequestURI),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
