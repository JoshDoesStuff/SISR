package api

import "net/http"

type StatusWriter struct {
	http.ResponseWriter
	Status int
	Error  error
}

func (w *StatusWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}
