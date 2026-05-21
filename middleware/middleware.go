package middleware

import "net/http"

func With(
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	// reverse loop to apply in the correct order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
