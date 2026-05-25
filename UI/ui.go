package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed build/*
var frontendFS embed.FS

func Middleware(preserveDocs bool, preserveOpenAPISpec bool) func(next http.Handler) http.Handler {
	subFS, err := fs.Sub(frontendFS, "build")

	if err != nil {
		panic(fmt.Sprintf("Build directory not found in frontendFS: %v", err))
	}

	s := http.FileServerFS(subFS)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}
			if preserveDocs && (r.URL.Path == "/docs" || r.URL.Path == "/docs/" || strings.HasPrefix(r.URL.Path, ("/schemas/"))) {
				next.ServeHTTP(w, r)
				return
			}
			if preserveOpenAPISpec {
				s := strings.Split(r.URL.Path, "/")
				if len(s) > 0 && strings.HasPrefix(s[len(s)-1], "openapi.") {
					next.ServeHTTP(w, r)
					return
				}
			}

			p := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
			_, err := subFS.Open(p)
			if err != nil {
				r.URL.Path = "/"
			}
			s.ServeHTTP(w, r)
		})
	}

}
