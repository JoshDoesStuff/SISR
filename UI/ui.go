package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
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
			serveIndex := p == "." || p == "" || p == "index.html"

			if !serveIndex {
				_, err := subFS.Open(p)
				if err != nil {
					serveIndex = true
				}
			}

			if serveIndex {
				apiPort := ""

				_, hostPort, err := net.SplitHostPort(r.Host)
				if err == nil {
					apiPort = hostPort
				} else if localAddr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
					_, localPort, localErr := net.SplitHostPort(localAddr.String())
					if localErr == nil {
						apiPort = localPort
					}
				}

				indexHTML, err := fs.ReadFile(subFS, "index.html")
				if err != nil {
					http.Error(w, "index.html not found", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(strings.ReplaceAll(string(indexHTML), "%SISR_API_PORT%", apiPort)))
				return
			}

			s.ServeHTTP(w, r)
		})
	}

}
