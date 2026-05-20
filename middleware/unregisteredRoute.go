package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Alia5/SISR/api"
	"github.com/danielgtaylor/huma/v2"
)

func UnregisteredRoute(a huma.API) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mux, ok := next.(*http.ServeMux)
		if ok {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, pattern := mux.Handler(r)
				if pattern == "" {
					pathExists := false
					for path, pathItem := range a.OpenAPI().Paths {
						if r.URL.Path == path && pathItem != nil {
							pathExists = true
							break
						}
					}
					var hErr huma.StatusError
					if pathExists {
						hErr = huma.Error405MethodNotAllowed("Method not allowed")
					} else {
						hErr = huma.Error404NotFound("Resource not found")
					}

					if sw, ok := w.(*api.StatusWriter); ok {
						sw.Error = hErr
					}

					b, err := json.Marshal(hErr)
					if err != nil {
						http.Error(w, "Internal server error", http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(hErr.GetStatus())
					_, _ = w.Write(b)
					return
				}
				mux.ServeHTTP(w, r)
			})
		}
		return next
	}
}
