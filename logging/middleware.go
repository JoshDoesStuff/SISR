package logging

import (
	"log/slog"
	"net/http"

	"github.com/Alia5/SISR/api/ctx"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		c := r.Context()
		statusCode := c.Value(ctx.KeyStatusCode)
		duration := c.Value(ctx.KeyDuration)

		err := c.Value(ctx.KeyError)
		if err != nil {
			slog.Error("request error",
				string(ctx.KeyStatusCode), statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", duration,
				"request_id", w.Header().Get("X-Request-ID"),
				"remote_addr", r.RemoteAddr,
				"error", err,
			)
			return
		}

		slog.Debug("request",
			string(ctx.KeyStatusCode), statusCode,
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", duration,
			"request_id", w.Header().Get("X-Request-ID"),
			"remote_addr", r.RemoteAddr,
		)
	})
}
